# train.py
import numpy as np
import os
import tensorflow as tf
from tensorflow.keras import layers, models
from matrix_builder import build_feature_matrix

def train_anomaly_detector():
    # Stream structural matrix arrays from memory builder logic
    X_train = build_feature_matrix(batch_size=100000)
    input_dim = X_train.shape[1] # Will be exactly 110 features
    
    print(f"Initializing Autoencoder Stack. Input Features: {input_dim}")
    
    # Hourglass structural architecture definitions
    model = models.Sequential([
        # Encoder Module: Compresses input patterns down
        layers.Input(shape=(input_dim,)),
        layers.Dense(64, activation='relu'),
        layers.Dense(32, activation='relu'),
        layers.Dense(8, activation='relu'),   # The Core Compressed Bottleneck Matrix
        
        # Decoder Module: Expands properties back out
        layers.Dense(32, activation='relu'),
        layers.Dense(64, activation='relu'),
        layers.Dense(input_dim, activation='sigmoid') # Reconstructs input dimensions (110)
    ])
    
    model.compile(optimizer='adam', loss='mse')
    model.summary()
    
    print("Launching Model Optimization Loops...")
    model.fit(
        X_train, X_train,
        epochs=5, 
        batch_size=1024, # Large batch sizes keep memory footprints stable on your Mac
        validation_split=0.1,
        shuffle=True
    )
    
    print("Evaluating system deviations to calculate operational threshold margins...")
    reconstructions = model.predict(X_train, batch_size=2048)
    mse_losses = np.mean(np.power(X_train - reconstructions, 2), axis=1)
    
    production_threshold = np.percentile(mse_losses, 98) 
    
    print(f"\n==========================================")
    print(f"98th Percentile Normal Deviation Loss: {production_threshold:.6f}")
    print(f"==========================================\n")

    target_dir = "model" 
    model_name = "vigilis_autoencoder.keras"

    os.makedirs(target_dir, exist_ok=True) 

    model_full_path = os.path.join(target_dir, model_name)
    threshold_full_path = os.path.join(target_dir, "threshold.npy")
    
    model.save(model_full_path)
    np.save(threshold_full_path, production_threshold)

    print("Successfully logged engine metadata parameters safely to localized workspace files.")

if __name__ == "__main__":
    train_anomaly_detector()