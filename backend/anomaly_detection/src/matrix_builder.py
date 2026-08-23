import numpy as np
from sklearn.feature_extraction.text import HashingVectorizer
from src.extractor import stream_logs_from_db

def get_one_hot_mapping():
    """
    Creates a static one-hot vector mapping for the 9 distinct callers.
    Each caller gets an identical, unbiased binary array representation.
    """
    distinct_callers = [
        "dfs.DataBlockScanner",
        "dfs.DataNode",
        "dfs.DataNode$BlockReceiver",
        "dfs.DataNode$DataTransfer",
        "dfs.DataNode$DataXceiver",
        "dfs.DataNode$PacketResponder",
        "dfs.FSDataset",
        "dfs.FSNamesystem",
        "dfs.PendingReplicationBlocks$PendingReplicationMonitor"
    ]
    
    mapping = {}
    # eg: mapping["dfs.DataBlockScanner"] = [1.0, 0., 0., 0., 0., 0., 0., 0., 0.] (np array)
    num_classes = len(distinct_callers) # Exactly 9 columns

    # one-hot encoding 
    for index, caller in enumerate(distinct_callers):
        one_hot = np.zeros(num_classes)
        one_hot[index] = 1.0
        mapping[caller] = one_hot
        
    return mapping

def build_feature_matrix(batch_size=100000):
    print("Building Feature Matrix from Live DB Stream...")
    
    vectorizer = HashingVectorizer(n_features=100, alternate_sign=False)
    one_hot_map = get_one_hot_mapping()
    fallback_vector = np.zeros(9) # Used if an unexpected caller appears
    
    all_matrices = []
    prev_timestamp = None
    
    for batch in stream_logs_from_db(batch_size=batch_size):
        if not batch:
            continue
            
        messages = [log['message'] for log in batch]
        
        # 1. Continuous Text Vectorization (100 columns)
        text_features = vectorizer.transform(messages).toarray()
        
        # 2. Extract Timing Context and One-Hot Encoded Categories
        meta_features = []
        for log in batch:
            # Fetch the unbiased 9-column binary representation
            caller_encoded = one_hot_map.get(log['caller'], fallback_vector)
            
            # Extract tracking intervals between sequential execution logs
            current_ts = log['ts']
            if prev_timestamp is None:
                time_delta = 0.0
            else:
                time_delta = (current_ts - prev_timestamp).total_seconds() * 1000.0
            
            prev_timestamp = current_ts
            normalized_time = np.log1p(max(0.0, time_delta)) / 20.0 
            
            # Glue the 9 binary columns and the 1 time column together (10 columns total)
            metadata_row = np.append(caller_encoded, normalized_time)
            meta_features.append(metadata_row)
            
        meta_features = np.array(meta_features)
        
        # 3. Final Stack: 100 text columns + 9 caller columns + 1 time column = 110 total features
        batch_matrix = np.hstack((text_features, meta_features))
        all_matrices.append(batch_matrix)
        print(f"Processed chunk size: {len(batch_matrix)} rows...")
        
    print("Stacking sub-matrices into full system array...")
    final_matrix = np.vstack(all_matrices)
    print(f"--- Matrix Allocation Finished --- Shape: {final_matrix.shape}")
    return final_matrix

if __name__ == "__main__":
    X = build_feature_matrix(batch_size=50000)