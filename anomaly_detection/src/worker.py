import re
import os
import json
import redis
import time
from datetime import datetime

import numpy as np
import tensorflow as tf
from sklearn.feature_extraction.text import HashingVectorizer

from prometheus_client import start_http_server,Counter, Histogram

from src.matrix_builder import get_one_hot_mapping
from config import config

def clean_log_message(message):

    message = re.sub(r'blk_-?\d+', 'BLOCK_ID', message)
    message = re.sub(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', 'IP_ADDR', message)
    message = re.sub(r':\d+', ':PORT', message)
    message = re.sub(r'\b\d+\b', 'NUM', message)
    return message

def initialize_ai_engine():
    print("Initializing Viligis Python Worker for anomaly_detection Autoencoder")
    
    model_path = config.MODEL_PATH
    threshold_path = config.THRESHOLD_PATH
    
    print(f"Loading Keras Model from {model_path}...") #!dev purpose, remove this
    model = tf.keras.models.load_model(model_path)
    
    print(f"Loading Anomaly Threshold from {threshold_path}...")#!dev purpose, remove this
    threshold = float(np.load(threshold_path))
    print(f"Loaded Threshold Value: {threshold:.6f}")

    # init vectorizer and the one hot map
    vectorizer = HashingVectorizer(n_features=100, alternate_sign=False)
    one_hot_map = get_one_hot_mapping()

    # connecting to the redis mq
    redis_host = config.REDIS_HOST
    redis_port = config.REDIS_PORT
    queue_name = config.REDIS_MQ_KEY_NAME

    #! redis init here
    r = redis.Redis(host=redis_host, port=redis_port, db=0, decode_responses=True)

    try:
        r.ping()
        print(f"Connected to Redis at {redis_host}:{redis_port}. Listening on key: '{queue_name}'")
    except redis.ConnectionError:
        raise RuntimeError("Failed to connect to Redis. Ensure Redis server is running!")

    vigilis_ai_total_logs_evaluated = Counter("vigilis_ai_total_logs_evaluated", "Total number of logs evaluated by the python autoencoder")
    vigilis_ai_total_anomalies_detected= Counter("vigilis_ai_total_anomalies_detected", "Number of anomaly logs found")
    vigilis_ai_inference_latency_seconds = Histogram("vigilis_ai_inference_latency_seconds", "Time it took (in seconds) for logs ingestion")

    return {
        "model": model,
        "threshold": threshold,
        "vectorizer": vectorizer,
        "one_hot_map": one_hot_map,

        "redis": r,
        "queue_name": queue_name,

        "vigilis_ai_total_logs_evaluated" : vigilis_ai_total_logs_evaluated,
        "vigilis_ai_total_anomalies_detected" : vigilis_ai_total_anomalies_detected,
        "vigilis_ai_inference_latency_seconds" : vigilis_ai_inference_latency_seconds,
    } #return the engine of objects

def build_live_batch_matrix(
    batch, vectorizer, one_hot_map, prev_timestamp=None
):
    # Transforms a live batch of JSON log dicts from Redis into an (N, 110)
    # feature matrix without touching matrix_builder.py.
    
    if not batch:
        return np.empty((0, 110)), prev_timestamp

    fallback_vector = np.zeros(9)

    cleaned_messages = [clean_log_message(log.get("msg", "")) for log in batch]
    text_features = vectorizer.transform(cleaned_messages).toarray()

    meta_features = []
    current_prev_ts = prev_timestamp

    for log in batch:
        caller_encoded = one_hot_map.get(
            log.get("caller", ""), fallback_vector
        )

        raw_ts = log.get("ts", "")
        try:
            if raw_ts.endswith("Z"):
                raw_ts = raw_ts[:-1]
            current_ts = datetime.fromisoformat(raw_ts)

            if current_prev_ts is None:
                time_delta = 0.0
            else:
                time_delta = (
                    current_ts - current_prev_ts
                ).total_seconds() * 1000.0

            current_prev_ts = current_ts
        except (ValueError, TypeError):
            time_delta = 0.0

        normalized_time = np.log1p(max(0.0, time_delta)) / 20.0

        metadata_row = np.append(caller_encoded, normalized_time)
        meta_features.append(metadata_row)

    meta_features = np.array(meta_features)

    # 100 text + 9 caller + 1 time = 110 total features
    batch_matrix = np.hstack((text_features, meta_features))

    return batch_matrix, current_prev_ts

def process_predictions(X_batch, batch, engine):
    r = engine["redis"]
    model = engine["model"]
    threshold = engine["threshold"]

    vigilis_ai_total_logs_evaluated = engine["vigilis_ai_total_logs_evaluated"]
    vigilis_ai_total_anomalies_detected = engine["vigilis_ai_total_anomalies_detected"]
    vigilis_ai_inference_latency_seconds = engine["vigilis_ai_inference_latency_seconds"]

    vigilis_ai_total_logs_evaluated.inc(len(batch))

    if X_batch.shape[0] == 0:
        return
    
    with vigilis_ai_inference_latency_seconds.time():
        reconstructed = model.predict(X_batch, verbose=0)

    mse_losses = np.mean(np.square(X_batch - reconstructed), axis=1)

    anomalies_found = 0
    for loss, log in zip(mse_losses, batch):
        if loss > threshold:
            anomalies_found += 1
            print("\n[SECURITY ALERT] Anomaly Detected!")
            print(f"   ├─ Timestamp: {log.get('ts')}")
            print(f"   ├─ Caller:    {log.get('caller')}")
            print(f"   ├─ Message:   {log.get('msg')}")
            print(
                f"   └─ MSE Loss:  {loss:.6f} (Threshold: {threshold:.6f})\n"
            )
            payload = {
                "level": log.get('level'),
                "timestamp": log.get('ts'),
                "caller": log.get('caller'),
                "message":log.get('msg'),
                "mse": f"{loss:.6f}",
                "threshold": f"{threshold:.6f}"
            }
            r.xadd(
                "log:anomaly",
                payload,
            )

    vigilis_ai_total_anomalies_detected.inc(anomalies_found)

    if anomalies_found == 0:
        print(
            f" clean batch {len(batch)} logs evaluated. Max MSE:"
            f" {np.max(mse_losses):.6f}"
        )

def start_worker_loop():
    engine = initialize_ai_engine()
    start_http_server(8000)
    r = engine["redis"]
    queue_name = engine["queue_name"]

    print("\n🟢🟢 Vigilis Autoencoder Ingestion Engine is live!!!\n")

    prev_timestamp = None

    while True:
        try:
            _, raw_payload = r.brpop(queue_name, timeout=0)
            batch = json.loads(raw_payload) # add .decode("utf-8") when using decode_responses=True parameter

            X_batch, prev_timestamp = build_live_batch_matrix(
                batch,
                engine["vectorizer"],
                engine["one_hot_map"],
                prev_timestamp,
            )

            process_predictions(X_batch, batch, engine)

        except json.JSONDecodeError as e:
            print(f"Error: Failed to parse JSON payload from Redis: {e}")
        except KeyboardInterrupt:
            print("\nShutting down worker process(KeyboardInterrupt)...")
            break
        except Exception as e:
            print(f"Ingestion loop error: {e}")
            time.sleep(1) # Prevent tight error looping

if __name__ == "__main__":
    start_worker_loop()