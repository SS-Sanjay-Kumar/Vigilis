import re
import os
import psycopg2
from psycopg2.extras import DictCursor
from dotenv import load_dotenv

load_dotenv()

def clean_log_message(message):

    message = re.sub(r'blk_-?\d+', 'BLOCK_ID', message)
    message = re.sub(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', 'IP_ADDR', message)
    message = re.sub(r':\d+', ':PORT', message)
    message = re.sub(r'\b\d+\b', 'NUM', message)
    return message

def stream_logs_from_db(batch_size=50000):

    conn = psycopg2.connect(
        dbname=os.getenv("DB_NAME"), 
        user=os.getenv("DB_USER"), 
        password=os.getenv("DB_PASSWORD"), 
        host=os.getenv("DB_HOST"), 
        port=os.getenv("DB_PORT")
    )
    
    cursor = conn.cursor('sentinel_log_stream_cursor', cursor_factory=DictCursor)
    
    # here extracting only the "normal" data, so a warning or an error would have an higher reconst error
    query = "SELECT ts, caller, message FROM logs WHERE level = 'info'"
    cursor.execute(query)
    
    while True:
        rows = cursor.fetchmany(batch_size)
        if not rows:
            break
            
        processed_batch = []
        for row in rows:
            cleaned_msg = f"info {clean_log_message(row['message'])}"
            processed_batch.append({
                'ts': row['ts'],
                'caller': row['caller'],
                'message': cleaned_msg
            })
        yield processed_batch

    cursor.close()
    conn.close()

if __name__ == "__main__":
    print("testing extractor before letting train.py use the functions")
    stream = stream_logs_from_db(batch_size=5)
    sample_batch = next(stream)
    for entry in sample_batch:
        print(f"[{entry['caller']}] -> {entry['message']}")