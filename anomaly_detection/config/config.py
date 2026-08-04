import os
from dotenv import load_dotenv
load_dotenv()

DATABASE_URL = os.getenv(
    "DATABASE_URL", 
    f"postgresql://{os.getenv('DB_USER')}:{os.getenv('DB_PASSWORD')}@{os.getenv('DB_HOST', 'localhost')}:{os.getenv('DB_PORT', '5432')}/{os.getenv('DB_NAME')}"
)

REDIS_HOST = os.getenv("REDIS_HOST")
REDIS_PORT = os.getenv("REDIS_PORT")
REDIS_MQ_KEY_NAME = os.getenv("REDIS_MQ_KEY_NAME")

MODEL_PATH =  os.getenv("MODEL_PATH")
THRESHOLD_PATH =  os.getenv("THRESHOLD_PATH")