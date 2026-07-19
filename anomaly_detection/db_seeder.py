import os
from datetime import datetime, timezone
from datasets import load_dataset
import psycopg2
from psycopg2.extras import execute_values
from dotenv import load_dotenv

load_dotenv()

def seed_database_from_hugging_face():
    conn = psycopg2.connect(
        dbname=os.getenv("DB_NAME"), 
        user=os.getenv("DB_USER"), 
        password=os.getenv("DB_PASSWORD"), 
        host=os.getenv("DB_HOST"), 
        port=os.getenv("DB_PORT")
    )

    print("test dotenv ->",os.getenv("DB_NAME")) #! test print
    cursor = conn.cursor()
    
    print("info: connecting to hugging face dataset stream")
    dataset = load_dataset("logfit-project/HDFS_v1", split="train", streaming=True) #* streamin = True is important or else machine will go OOM
    
    batch_size = 50000
    batch = []
    total_inserted = 0
    print("starting execution stream and database insertion")
    
    for row in dataset:

        raw_log_time = datetime.strptime(f"{row['date']} {row['time']}", "%y%m%d %H%M%S").replace(tzinfo=timezone.utc)
        log_timestamp = raw_log_time.isoformat()

        log_entry = (
            row['level'].lower(),        
            log_timestamp,
            row['component'],           # caller
            row['content']              # message
        )
        batch.append(log_entry)
        
        if len(batch) >= batch_size:
            query = "INSERT INTO logs (level, ts, caller, message) VALUES %s"
            execute_values(cursor, query, batch)
            conn.commit()
            
            total_inserted += len(batch)
            print(f"Successfully committed {total_inserted} records to Postgres...")
            batch = []
            
    if batch:
        query = "INSERT INTO logs (timestamp, level, component, message) VALUES %s"
        execute_values(cursor, query, batch)
        conn.commit()
        total_inserted += len(batch)
        print(f"Final flush completed. Total rows seeded: {total_inserted}")
        
    cursor.close()
    conn.close()

if __name__ == "__main__":
    seed_database_from_hugging_face()