# evaluate_test_set.py
import os
import re
import numpy as np
import psycopg2
from psycopg2.extras import DictCursor
from sklearn.feature_extraction.text import HashingVectorizer
import tensorflow as tf
from src.matrix_builder import get_one_hot_mapping
from config import config


def clean_log_message(message):
    message = re.sub(r'blk_-?\d+', 'BLOCK_ID', message)
    message = re.sub(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', 'IP_ADDR', message)
    message = re.sub(r':\d+', ':PORT', message)
    message = re.sub(r'\b\d+\b', 'NUM', message)
    return message

def evaluate_engine():
    print("Loading models and threshold metrics from disk...")
    if not os.path.exists("model/vigilis_autoencoder.keras") or not os.path.exists("model/threshold.npy"):
        print("Model files missing! Please run train.py first.")
        return
        
    model = tf.keras.models.load_model("model/vigilis_autoencoder.keras")
    production_threshold = float(np.load("model/threshold.npy"))

    vectorizer = HashingVectorizer(n_features=100, alternate_sign=False)
    one_hot_map = get_one_hot_mapping()
    fallback_vector = np.zeros(9)
    
    total_processed = 0
    detected_anomalies = 0
    critical_alerts = 0
    warning_alerts = 0
    
    all_mse_losses = []
    prev_timestamp = None
    
    # Open the audit file to record anomalies as they stream in
    audit_filename = "alert_audit_report.md"
    
    print(f"Initializing audit log target file: {audit_filename}")
    with open(audit_filename, "w", encoding="utf-8") as report:
        report.write("# Sentinel-Stream Anomaly Detection Audit Report\n")
        report.write(f"**Baseline Production Threshold:** `{production_threshold:.6f}`\n\n")
        report.write("## 🚨 CRITICAL ALERTS (Risk >= 75%)\n")
        report.write("| Timestamp | Level | Caller | MSE | Risk Factor | Raw Message |\n")
        report.write("| --- | --- | --- | --- | --- | --- |\n")
    
    print("Connecting to database and initializing streaming pipeline...")
    conn = psycopg2.connect(config.DATABASE_URL)
    
    cursor = conn.cursor(name="anomaly_stream_cursor", cursor_factory=DictCursor)
    query = "SELECT ts, caller, level, message FROM logs WHERE level != 'info'"
    cursor.execute(query)
    
    chunk_size = 50000
    cursor.itersize = chunk_size 
    
    # Lists to hold standard warning rows during streaming so we don't mess up Markdown table sequence
    warning_buffer = []
    
    print(f"Beginning continuous evaluation processing in streams of {chunk_size} records...")
    
    while True:
        rows = cursor.fetchmany(chunk_size)
        if not rows:
            break
            
        messages = [f"{clean_log_message(row['message'])}" for row in rows]
        text_features = vectorizer.transform(messages).toarray()
        
        meta_features = []
        for row in rows:
            caller_encoded = one_hot_map.get(row['caller'], fallback_vector)
            
            current_ts = row['ts']
            if prev_timestamp is None:
                time_delta = 0.0
            else:
                time_delta = (current_ts - prev_timestamp).total_seconds() * 1000.0
            
            prev_timestamp = current_ts
            normalized_time = np.log1p(max(0.0, time_delta)) / 20.0 
            
            metadata_row = np.append(caller_encoded, normalized_time)
            meta_features.append(metadata_row)
            
        meta_features = np.array(meta_features)
        X_chunk = np.hstack((text_features, meta_features))
        
        reconstructions = model.predict(X_chunk, batch_size=1024, verbose=0)
        mse_losses = np.mean(np.power(X_chunk - reconstructions, 2), axis=1)
        
        all_mse_losses.extend(mse_losses)
        
        MSE_MAX_SCALE = 0.015000 
        
        # Open in append mode to log rows from this chunk instantly
        with open(audit_filename, "a", encoding="utf-8") as report:
            for idx, row in enumerate(rows):
                mse = mse_losses[idx]
                level = row['level'].upper()
                
                if mse > production_threshold:
                    detected_anomalies += 1
                    
                    raw_score = ((mse - production_threshold) / (MSE_MAX_SCALE - production_threshold)) * 100.0
                    raw_score = max(0.0, raw_score)
                    
                    if level in ['ERROR', 'FATAL']:
                        severity_multiplier = 1.0
                    elif level == 'WARN':
                        severity_multiplier = 0.5
                    else:
                        severity_multiplier = 0.2
                        
                    risk_percentage = min(100.0, raw_score * severity_multiplier)
                    
                    clean_msg = row['message'].replace("\n", " ").replace("|", "\\|")
                    
                    if risk_percentage >= 75.0:
                        critical_alerts += 1
                        report.write(f"| {row['ts']} | {level} | `{row['caller']}` | {mse:.6f} | {risk_percentage:.2f}% | {clean_msg} |\n")
                    elif risk_percentage >= 40.0:
                        warning_alerts += 1
                        # Buffer warnings to print them cleanly in their own section later
                        warning_buffer.append(f"| {row['ts']} | {level} | `{row['caller']}` | {mse:.6f} | {risk_percentage:.2f}% | {clean_msg} |\n")
                        
        total_processed += len(rows)
        print(f"Processed: {total_processed} vectors...")

    # Write out buffered Warning logs to the end of the report file
    with open(audit_filename, "a", encoding="utf-8") as report:
        report.write("\n\n## ⚠️ WARNING ALERTS (Risk 40%-74%)\n")
        report.write("| Timestamp | Level | Caller | MSE | Risk Factor | Raw Message |\n")
        report.write("| --- | --- | --- | --- | --- | --- |\n")
        for line in warning_buffer:
            report.write(line)

    cursor.close()
    conn.close()
    
    all_mse_losses = np.array(all_mse_losses)
    detection_rate = (detected_anomalies / total_processed) * 100.0
    
    print(f"\n================ FINAL PRODUCTION SYSTEM METRICS ================")
    print(f"Configured Safe System Threshold       : {production_threshold:.6f}")
    print(f"Total Threat Vectors Processed         : {total_processed}")
    print(f"Average Anomaly Reconstruction Error  : {np.mean(all_mse_losses):.6f}")
    print(f"Maximum Anomaly Reconstruction Error  : {np.max(all_mse_losses):.6f}")
    print(f"Successfully Flagged Vectors          : {detected_anomalies} / {total_processed} records")
    print(f"Targeted System Detection Accuracy    : {detection_rate:.2f}%")
    print(f"----------------------------------------------------------------")
    print(f"CRITICAL ALERTS ROUTED (Risk >= 75%)   : {critical_alerts}")
    print(f"WARNING ALERTS ROUTED  (Risk 40%-74%)  : {warning_alerts}")
    print(f"=================================================================\n")
    print(f"Audit generation complete. Please review: {audit_filename}")

if __name__ == "__main__":
    evaluate_engine()