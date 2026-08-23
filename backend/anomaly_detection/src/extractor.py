import re
import os
import psycopg2
from psycopg2.extras import DictCursor
from config import config

def clean_log_message(message):

    message = re.sub(r'blk_-?\d+', 'BLOCK_ID', message)
    message = re.sub(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', 'IP_ADDR', message)
    message = re.sub(r':\d+', ':PORT', message)
    message = re.sub(r'\b\d+\b', 'NUM', message)
    return message

def stream_logs_from_db(batch_size=50000):

    conn = psycopg2.connect(config.DATABASE_URL)
    
    cursor = conn.cursor('vigilis_log_stream_cursor', cursor_factory=DictCursor)
    # named cursor = server side cursor
    # server side cursor: the server runs the SQL command and holds the data(generally large) in its memory buffer
    # it doesnt does the data to the client now
    # we can use the cursor to fetch the data in batches
    # anolagy: cursor(named) asks postgres -> Hey, remember that specific query I ran a minute ago? Give me the next 100 rows from it.
    # this is used when data to fetch is significantly large and may result in OOM

    # cursor_factory =it is a configuration setting in the database driver that changes the data structure used to return the query results.
    
    # here extracting only the "normal" data, so a warning or an error would have an higher reconst error
    query = "SELECT ts, caller, message FROM logs WHERE level = 'info'"
    cursor.execute(query)
    
    while True:
        rows = cursor.fetchmany(batch_size) #fetching in batches, only can be done with named cursors
        if not rows:
            break
            
        processed_batch = []
        for row in rows:
            cleaned_msg = f"{clean_log_message(row['message'])}"
            processed_batch.append({
                'ts': row['ts'],
                'caller': row['caller'],
                'message': cleaned_msg,
            })
        yield processed_batch

    cursor.close()
    conn.close()

if __name__ == "__main__": #only for testing purposes, only runs when this file is ran individually
    print("testing extractor before letting train.py use the functions")

    stream = stream_logs_from_db(batch_size=5)  
    sample_batch = next(stream) #cause we used yield
    for entry in sample_batch:
        print(f"[{entry['caller']}] -> {entry['message']}")