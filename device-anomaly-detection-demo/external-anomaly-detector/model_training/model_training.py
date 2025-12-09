import bnlearn as bn
import pandas as pd
import os

# Define paths
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DATA_PATH = os.path.join(BASE_DIR, 'data', 'observed_data.csv')
MODEL_DIR = os.path.join(BASE_DIR, 'models')
MODEL_PATH = os.path.join(MODEL_DIR, 'bn_model.pkl')

def model_training():
    """
    Train a Bayesian Network model using data from observed_data.csv.
    Performs structure learning and parameter learning.
    Saves the trained model to the models directory.
    """
    print(f"Loading data from {DATA_PATH}...")
    df = pd.read_csv(DATA_PATH)
    
    # Structure learning
    print("Performing structure learning...")
    # methodtype='hc' (Hill Climbing) or 'ex' (ExhaustiveSearch) or 'cs' (Constraint-based)
    # 'hc' is a good default for structure learning
    dag = bn.structure_learning.fit(df, methodtype='hc', verbose=0)
    
    # Parameter learning
    print("Performing parameter learning...")
    model = bn.parameter_learning.fit(dag, df, verbose=0)
    
    # Save model
    if not os.path.exists(MODEL_DIR):
        os.makedirs(MODEL_DIR)
    
    print(f"Saving model to {MODEL_PATH}...")
    bn.save(model, filepath=MODEL_PATH, overwrite=True)
    print("Model training completed.")
    return model

def predict_brightness(light1, light2):
    """
    Predict the value of brightness given light1 and light2 values using the trained model.
    """
    # Load model
    if os.path.exists(MODEL_PATH):
        model = bn.load(filepath=MODEL_PATH, verbose=0)
    else:
        print("Model not found. Training new model...")
        model = model_training()
        
    # Inference
    # Calculate P(brightness | light1, light2)
    query = bn.inference.fit(model, variables=['brightness'], evidence={'light1': light1, 'light2': light2}, verbose=0)
    
    # Extract the most likely value for brightness
    if hasattr(query, 'df'):
        # query.df contains columns for the variable (brightness) and 'p' (probability)
        # We want the value of brightness with the highest probability
        max_prob_row = query.df.loc[query.df['p'].idxmax()]
        predicted_brightness = max_prob_row['brightness']
        return predicted_brightness
    else:
        print("Error: Could not retrieve prediction results.")
        return None

if __name__ == "__main__":
    model_training()
    
    # Simple test
    print("\nTesting prediction:")
    for l1 in [0, 1]:
        for l2 in [0, 1]:
            pred = predict_brightness(l1, l2)
            print(f"Light1={l1}, Light2={l2} -> Predicted Brightness={pred}")
