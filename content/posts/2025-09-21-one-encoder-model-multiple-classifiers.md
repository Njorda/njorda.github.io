---
layout: post 
title: "How to Train a Text Classifier Without Having an Embedding Model or GPU Locally"
subtitle: "A Step-by-Step Guide to Building a Text Classifier Using an Embedding API"
date: 2025-09-21
author: "Johan Hansson"
URL: "/2025/09/21/one-encoder-model-multiple-classifiers/"
image: "/img/background_2022_07_17.png"
---


In this example, we will go over how to build a text classifier using a text embedding API. This article is inspired by [Fine-tune classifier with ModernBERT in 2025](https://www.philschmid.de/fine-tune-modern-bert-in-2025)

```
Large Language Models (LLMs) have become ubiquitous in 2024. However, smaller, specialized models - particularly for classification tasks - remain critical for building efficient and cost-effective AI systems. One key use case is routing user prompts to the most appropriate LLM or selecting optimal few-shot examples, where fast, accurate classification is essential.
```

But instead of running an encoder model locally, we will utilize an API to get embeddings and then train a small classifier. This will create a very lightweight setup where you can build multiple classifiers on top of one encoder.

![One encoder model multiple classifiers](/img/embedding_architecture_diagram.svg)

The diagram above illustrates the key advantage of our approach:

**Traditional Setup:** Each classifier requires its own fine-tuned model running on expensive GPU infrastructure

**Our Approach:** One centralized encoder model provides embeddings via API to multiple lightweight CPU-based classifiers

This architecture provides several benefits:

💡 **Lower Costs:** Reduced GPU infrastructure requirements  
🎯 **Centralized Management:** Single point of maintenance for the encoder  
🔧 **Easy Scaling:** Add new classifiers without additional GPU resources

While this API-based approach offers significant advantages, there are some important considerations:

- **Latency:** API calls add network delay compared to local inference
- **Embedding Model Lock-in:** Changing the embedding model requires retraining all classifiers
- **Dependency:** Requires stable internet and API availability
- **No Fine-tuning Control:** Can't optimize the encoder for your specific domain/task
- **Version Dependencies:** API updates to the embedding model could break existing classifiers

Okay, with that said, let's start building! We will use the same dataset as the article linked above: [Banking77Classification](https://huggingface.co/datasets/mteb/banking77). It has what we need—a text and a label. The first step is to download the dataset and add embeddings to the dataset, one embedding for each text. In this example, we will use Gemini (GCP) to create the embeddings, but you can use any embedding API you want. However, remember that different embeddings will give different results, and that might affect your end performance.

```python 


import numpy as np
import time
import os
from google import genai
from google.genai import types
from tqdm.auto import tqdm  # Use tqdm.auto for better Jupyter support
from datasets import Dataset, DatasetDict
from datasets import load_dataset
from dotenv import load_dotenv


load_dotenv()
api_key = os.getenv('gemini_api_key')

# Dataset id from huggingface.co/dataset
# dataset_id = "DevQuasar/llm_router_dataset-synth"
dataset_id = "legacy-datasets/banking77"

if not api_key:
    raise ValueError("gemini_api_key not found in environment variables")


client = genai.Client(api_key=api_key)

def get_embeddings_batch(texts, max_retries=3):
    """
    Get embeddings for a batch of texts with retry logic and exponential backoff.

    Args:
        texts: List of text strings to embed
        max_retries: Maximum number of retry attempts

    Returns:
        List of numpy arrays containing embeddings
    """
    for attempt in range(max_retries):
        try:
            result = client.models.embed_content(
                model="gemini-embedding-001",
                contents=texts
            )
            return [np.array(e.values) for e in result.embeddings]
        except Exception as e:
            print(f"Attempt {attempt + 1} failed: {str(e)}")
            
            if attempt == max_retries - 1:
                # Last attempt failed, re-raise the exception
                raise e
            
            # Exponential backoff: wait 2^attempt seconds
            wait_time = 2 ** attempt
            print(f"Waiting {wait_time} seconds before retry...")
            time.sleep(wait_time)

def add_embeddings_to_dataset(dataset, text_column="text", batch_size=50, test_mode=False, test_samples=10):
    """
    Add embeddings to a Hugging Face dataset with improved rate limiting.

    Args:
        dataset: Hugging Face dataset object
        text_column: Name of the column containing text to embed
        batch_size: Number of texts to process at once (reduced to 50 for better stability)
        test_mode: If True, only process a small subset for testing
        test_samples: Number of samples to process in test mode

    Returns:
        Dataset with added 'embedding' column
    """

    def process_split(split_data, split_name):
        """Process a single dataset split (train/test)."""
        print(f"\nProcessing {split_name} split...")

        # Get texts to process
        if test_mode:
            texts = split_data[text_column][:test_samples]
            total_samples = min(test_samples, len(split_data))
        else:
            texts = split_data[text_column]
            total_samples = len(split_data)

        print(f"Total samples to process: {total_samples}")

        # Process in batches
        all_embeddings = []

        for i in tqdm(range(0, total_samples, batch_size), desc=f"Embedding {split_name}"):
            batch_end = min(i + batch_size, total_samples)
            batch_texts = texts[i:batch_end]

            try:
                # Get embeddings for this batch
                batch_embeddings = get_embeddings_batch(batch_texts)
                all_embeddings.extend(batch_embeddings)
                
                # Progress update
                print(f"Processed batch {i//batch_size + 1}/{(total_samples + batch_size - 1)//batch_size}")

            except Exception as e:
                print(f"Failed to process batch starting at index {i}: {str(e)}")
                # You could choose to skip this batch or stop processing
                # For now, we'll skip and continue
                continue

            # Rate limiting: wait between batches to avoid hitting API limits
            if i + batch_size < total_samples:
                time.sleep(1)  # Reduced from 0.5 to 1 second for better rate limiting

        # Add embeddings to the dataset
        if test_mode:
            # For test mode, create a subset with embeddings
            subset_data = {key: split_data[key][:len(all_embeddings)] for key in split_data.features}
            subset_data['embedding'] = all_embeddings
            return subset_data
        else:
            # Add embeddings to full dataset
            split_data = split_data.add_column('embedding', all_embeddings)
            return split_data

    # Process each split
    processed_dataset = {}

    for split_name in dataset.keys():
        processed_dataset[split_name] = process_split(dataset[split_name], split_name)

    return processed_dataset

def main():
    """Main function to load dataset and add embeddings."""

    # Configuration
    DATASET_ID = "legacy-datasets/banking77"
    TEXT_COLUMN = "text"
    BATCH_SIZE = 100  # Reduced from 100 to 50 for better stability
    TEST_MODE = False  # Set to False to process full dataset
    TEST_SAMPLES = 100  # Number of samples to process in test mode

    print("=" * 50)
    print("Dataset Embedding Generator")
    print("=" * 50)

    # Load dataset
    print(f"\nLoading dataset: {DATASET_ID}")
    raw_dataset = load_dataset(DATASET_ID)

    print(f"Train dataset size: {len(raw_dataset['train'])}")
    print(f"Test dataset size: {len(raw_dataset['test'])}")

    if TEST_MODE:
        print(f"\n⚠️  TEST MODE: Processing only {TEST_SAMPLES} samples per split")
        print("Set TEST_MODE = False to process the full dataset")

    # Add embeddings to dataset
    dataset_with_embeddings = add_embeddings_to_dataset(
        raw_dataset,
        text_column=TEXT_COLUMN,
        batch_size=BATCH_SIZE,
        test_mode=TEST_MODE,
        test_samples=TEST_SAMPLES
    )

    # Save the dataset with embeddings FIRST (to avoid losing work)
    print("\n" + "=" * 50)
    print("Saving dataset with embeddings...")
    print("=" * 50)

    # Convert to Hugging Face Dataset format if needed
    from datasets import Dataset, DatasetDict

    final_dataset = DatasetDict({
        split: Dataset.from_dict(data) if isinstance(data, dict) else data
        for split, data in dataset_with_embeddings.items()
    })

    # Save to disk
    final_dataset.save_to_disk("datasets/banking77_with_embeddings_test")
    print("Dataset saved to datasets/banking77_with_embeddings_test")

    # Verify embeddings were added (after saving)
    print("\n" + "=" * 50)
    print("Verification")
    print("=" * 50)

    for split_name in dataset_with_embeddings.keys():
        if TEST_MODE:
            # For test mode, check the dictionary structure
            sample_embedding = dataset_with_embeddings[split_name]['embedding'][0]
            num_samples = len(dataset_with_embeddings[split_name]['embedding'])
        else:
            # For full mode, check the dataset structure
            sample_embedding = dataset_with_embeddings[split_name][0]['embedding']
            num_samples = len(dataset_with_embeddings[split_name])

        # Convert to numpy array to get shape safely
        import numpy as np
        if isinstance(sample_embedding, list):
            sample_embedding = np.array(sample_embedding)

        print(f"\n{split_name.capitalize()} split:")
        print(f"  - Number of samples: {num_samples}")
        print(f"  - Embedding shape: {sample_embedding.shape}")
        print(f"  - Sample text: {dataset_with_embeddings[split_name][TEXT_COLUMN][0][:100]}...")
        print(f"  - Sample embedding (first 5 dims): {sample_embedding[:5]}")

    
    return dataset_with_embeddings

def convert_to_huggingface_dataset(dataset_with_embeddings):
    """
    Convert the processed dataset (from test or full mode) into a proper Hugging Face Dataset.

    Args:
        dataset_with_embeddings: Dictionary with train/test splits containing embeddings

    Returns:
        DatasetDict: Properly formatted Hugging Face dataset
    """

    print("Converting to Hugging Face Dataset format...")

    # Create a new DatasetDict
    final_dataset_dict = {}

    for split_name, split_data in dataset_with_embeddings.items():
        print(f"Processing {split_name} split...")

        if isinstance(split_data, dict):
            # Test mode or dictionary format
            # Convert numpy arrays to lists for proper serialization
            data_dict = {}
            for key, values in split_data.items():
                if key == 'embedding':
                    # Convert numpy arrays to lists
                    data_dict[key] = [arr.tolist() if isinstance(arr, np.ndarray) else arr
                                     for arr in values]
                else:
                    data_dict[key] = values

            # Create Dataset from dictionary
            dataset = Dataset.from_dict(data_dict)

        else:
            # Already a Dataset object, but need to ensure embeddings are properly formatted
            # Convert the dataset to dict, process embeddings, then recreate
            data_dict = split_data.to_dict() if hasattr(split_data, 'to_dict') else dict(split_data)

            # Ensure embeddings are in list format
            if 'embedding' in data_dict:
                data_dict['embedding'] = [
                    arr.tolist() if isinstance(arr, np.ndarray) else arr
                    for arr in data_dict['embedding']
                ]

            dataset = Dataset.from_dict(data_dict)

        final_dataset_dict[split_name] = dataset

        # Print info about this split
        print(f"  - {split_name}: {len(dataset)} samples")
        if 'embedding' in dataset.column_names:
            sample_embedding = dataset[0]['embedding']
            print(f"  - Embedding dimension: {len(sample_embedding)}")

    # Create final DatasetDict
    final_dataset = DatasetDict(final_dataset_dict)

    print("\nDataset conversion complete!")
    print(f"Dataset info: {final_dataset}")

    return final_dataset


main()

```
This code will read the dataset and add embeddings to each text field. It will send batches of text to the API to improve latency and, in this case, also lower the cost (lower cost per token for batch requests). However, for Gemini and my quota, there is a maximum of 100 text snippets in each request. Finally, we store the new dataset on disk so we don't have to recreate it. If you want to optimize the code in the future, you can update it to run async—this would speed it up; however, make sure you have the quota/API limits to support it.

The next step is to build the classifier. There are multiple types of models that would work fine, but a good starting point is normally Random Forest. Below are some of its properties:

- **Minimal preprocessing required** - Unlike many algorithms, Random Forest handles mixed data types (numerical and categorical) well without extensive feature scaling or encoding. It's also relatively robust to outliers and missing values.
- **Built-in feature selection** - The algorithm provides feature importance scores automatically, helping you understand which variables matter most for your predictions. This insight is valuable for both model interpretation and feature engineering.
- **Reduces overfitting** - By combining multiple decision trees and using bootstrap sampling, Random Forest naturally reduces the overfitting that individual decision trees are prone to, often giving you better generalization without much tuning.
- **Good default performance** - Random Forest typically performs well "out of the box" with minimal hyperparameter tuning. You can often get reasonable results with default settings before optimizing further.
- **Handles non-linear relationships** - The tree-based structure captures complex, non-linear patterns in data without requiring you to manually specify interaction terms or transformations.
- **Provides uncertainty estimates** - The probability outputs from Random Forest give you a sense of prediction confidence, which is useful for understanding model reliability.
- **Computationally reasonable** - While not the fastest algorithm, Random Forest trains relatively quickly and can handle moderately large datasets without specialized hardware.
- **Interpretable foundation** - Though not as interpretable as a single decision tree, Random Forest is more explainable than black-box methods like neural networks, making it easier to validate and trust your initial results.

Of course, it has some drawbacks, where the biggest one is that it can struggle with highly imbalanced datasets. However, in our case, that is not a problem. Compared to other model types, it is slow at prediction time since we will have to traverse multiple trees, but for this case, it is fine.

Below is the code for training a Random Forest model. 
```python 

import numpy as np
import joblib
import os
from datasets import load_from_disk
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report, accuracy_score, confusion_matrix
from sklearn.model_selection import cross_val_score
import matplotlib.pyplot as plt
import seaborn as sns
from tqdm.auto import tqdm
import warnings
warnings.filterwarnings('ignore')

def load_dataset_with_embeddings(dataset_path="banking77_with_embeddings_test"):
    """
    Load the dataset with embeddings from disk.
    
    Args:
        dataset_path: Path to the saved dataset directory
        
    Returns:
        DatasetDict containing train and test splits with embeddings
    """
    print(f"Loading dataset from {dataset_path}...")
    dataset = load_from_disk(dataset_path)
    
    print(f"Train samples: {len(dataset['train'])}")
    print(f"Test samples: {len(dataset['test'])}")
    print(f"Number of classes: {len(dataset['train'].features['label'].names)}")
    print(f"Embedding dimension: {len(dataset['train'][0]['embedding'])}")
    
    return dataset

def prepare_data(dataset, test_mode=False, test_samples=100):
    """
    Prepare the data for training by extracting embeddings and labels.
    
    Args:
        dataset: DatasetDict with train and test splits
        test_mode: If True, use only a subset of data for testing
        test_samples: Number of samples to use in test mode
        
    Returns:
        Tuple of (X_train, y_train, X_test, y_test, label_names)
    """
    print("\nPreparing data...")
    
    # Extract train data
    if test_mode:
        train_data = dataset['train'].select(range(min(test_samples, len(dataset['train']))))
        test_data = dataset['test'].select(range(min(test_samples, len(dataset['test']))))
        print(f"Using {len(train_data)} train samples and {len(test_data)} test samples (test mode)")
    else:
        train_data = dataset['train']
        test_data = dataset['test']
        print(f"Using {len(train_data)} train samples and {len(test_data)} test samples (full mode)")
    
    # Extract embeddings and labels
    X_train = np.array([sample['embedding'] for sample in train_data])
    y_train = np.array([sample['label'] for sample in train_data])
    
    X_test = np.array([sample['embedding'] for sample in test_data])
    y_test = np.array([sample['label'] for sample in test_data])
    
    # Get label names
    label_names = train_data.features['label'].names
    
    print(f"Training data shape: {X_train.shape}")
    print(f"Test data shape: {X_test.shape}")
    print(f"Number of unique labels in train: {len(np.unique(y_train))}")
    print(f"Number of unique labels in test: {len(np.unique(y_test))}")
    
    return X_train, y_train, X_test, y_test, label_names

def train_random_forest(X_train, y_train, n_estimators=100, max_depth=None, random_state=42, test_mode=False):
    """
    Train a Random Forest classifier on the embeddings.
    
    Args:
        X_train: Training embeddings
        y_train: Training labels
        n_estimators: Number of trees in the forest
        max_depth: Maximum depth of trees
        random_state: Random state for reproducibility
        test_mode: If True, use smaller parameters for faster training
        
    Returns:
        Trained RandomForestClassifier
    """
    print(f"\nTraining Random Forest classifier...")
    
    if test_mode:
        # Use smaller parameters for test mode
        n_estimators = min(50, n_estimators)
        max_depth = min(10, max_depth) if max_depth else 10
        print(f"Test mode: Using {n_estimators} estimators, max_depth={max_depth}")
    
    # Create and train the classifier
    rf_classifier = RandomForestClassifier(
        n_estimators=n_estimators,
        max_depth=max_depth,
        random_state=random_state,
        n_jobs=-1,  # Use all available cores
        verbose=1 if test_mode else 0
    )
    
    print("Fitting the model...")
    rf_classifier.fit(X_train, y_train)
    
    print("Training completed!")
    return rf_classifier

def evaluate_model(model, X_test, y_test, label_names, test_mode=False):
    """
    Evaluate the trained model and print detailed results.
    
    Args:
        model: Trained RandomForestClassifier
        X_test: Test embeddings
        y_test: Test labels
        label_names: List of class names
        test_mode: If True, skip some detailed analysis for speed
    """
    print("\n" + "="*60)
    print("MODEL EVALUATION")
    print("="*60)
    
    # Make predictions
    print("Making predictions...")
    y_pred = model.predict(X_test)
    
    # Calculate accuracy
    accuracy = accuracy_score(y_test, y_pred)
    print(f"\nOverall Accuracy: {accuracy:.4f} ({accuracy*100:.2f}%)")
    
    # Get unique classes in test data and their corresponding names
    unique_classes = np.unique(np.concatenate([y_test, y_pred]))
    relevant_label_names = [label_names[i] for i in unique_classes]
    
    print(f"\nClasses present in test data: {len(unique_classes)}")
    print(f"Class indices: {unique_classes}")
    print(f"Class names: {relevant_label_names}")
    
    # Classification report
    print("\nDetailed Classification Report:")
    print("-" * 40)
    report = classification_report(
        y_test, y_pred, 
        target_names=relevant_label_names, 
        labels=unique_classes,
        zero_division=0
    )
    print(report)
    
    # Save classification report to file
    report_filename = "classification_report.txt"
    with open(report_filename, 'w') as f:
        f.write("Classification Report\n")
        f.write("=" * 50 + "\n")
        f.write(f"Overall Accuracy: {accuracy:.4f} ({accuracy*100:.2f}%)\n\n")
        f.write(f"Classes present in test data: {len(unique_classes)}\n")
        f.write(f"Class indices: {unique_classes}\n")
        f.write(f"Class names: {relevant_label_names}\n\n")
        f.write("Detailed Classification Report:\n")
        f.write("-" * 40 + "\n")
        f.write(report)
    print(f"Classification report saved to '{report_filename}'")
    
    # Feature importance (top 10)
    print("\nTop 10 Most Important Features (Embedding Dimensions):")
    print("-" * 50)
    feature_importance = model.feature_importances_
    top_indices = np.argsort(feature_importance)[-10:][::-1]
    
    for i, idx in enumerate(top_indices):
        print(f"{i+1:2d}. Dimension {idx:4d}: {feature_importance[idx]:.6f}")
    
    # Always create confusion matrix
    print("\nConfusion Matrix:")
    print("-" * 20)
    cm = confusion_matrix(y_test, y_pred, labels=unique_classes)
    print(cm)
    
    # Save confusion matrix as text file
    cm_filename = "confusion_matrix.txt"
    with open(cm_filename, 'w') as f:
        f.write("Confusion Matrix\n")
        f.write("=" * 30 + "\n")
        f.write(f"Shape: {cm.shape}\n")
        f.write(f"Classes: {len(unique_classes)}\n\n")
        f.write("Matrix (rows=actual, cols=predicted):\n")
        f.write(str(cm))
        f.write(f"\n\nClass indices: {unique_classes}\n")
        f.write(f"Class names: {relevant_label_names}\n")
    print(f"Confusion matrix saved to '{cm_filename}'")
    
    # Create confusion matrix visualization
    if len(unique_classes) <= 15:
        # For small number of classes, create detailed heatmap
        plt.figure(figsize=(12, 10))
        sns.heatmap(cm, annot=True, fmt='d', cmap='Blues', 
                   xticklabels=relevant_label_names, 
                   yticklabels=relevant_label_names)
        plt.title('Confusion Matrix')
        plt.ylabel('True Label')
        plt.xlabel('Predicted Label')
        plt.xticks(rotation=45, ha='right')
        plt.yticks(rotation=0)
        plt.tight_layout()
        plt.savefig('confusion_matrix.png', dpi=150, bbox_inches='tight')
        print("Confusion matrix visualization saved as 'confusion_matrix.png'")
        plt.close()
    else:
        # For large number of classes, create a summary heatmap
        plt.figure(figsize=(15, 12))
        sns.heatmap(cm, annot=False, fmt='d', cmap='Blues')
        plt.title(f'Confusion Matrix ({len(unique_classes)} classes)')
        plt.ylabel('True Label Index')
        plt.xlabel('Predicted Label Index')
        plt.tight_layout()
        plt.savefig('confusion_matrix.png', dpi=150, bbox_inches='tight')
        print("Confusion matrix visualization saved as 'confusion_matrix.png' (large dataset - no class labels)")
        plt.close()
    
    return accuracy, y_pred

def save_model(model, model_path="random_forest_model.joblib"):
    """
    Save the trained model to disk.
    
    Args:
        model: Trained RandomForestClassifier
        model_path: Path to save the model
    """
    print(f"\nSaving model to {model_path}...")
    joblib.dump(model, model_path)
    print("Model saved successfully!")

def load_model(model_path="random_forest_model.joblib"):
    """
    Load a trained model from disk.
    
    Args:
        model_path: Path to the saved model
        
    Returns:
        Loaded RandomForestClassifier
    """
    print(f"Loading model from {model_path}...")
    model = joblib.load(model_path)
    print("Model loaded successfully!")
    return model


def main():
    """
    Main function to run the complete pipeline.
    """
    print("="*60)
    print("RANDOM FOREST CLASSIFIER FOR BANKING77 WITH EMBEDDINGS")
    print("="*60)
    
    # Configuration
    TEST_MODE = False  # Set to False to use full dataset
    TEST_SAMPLES = 200  # Number of samples to use in test mode
    DATASET_PATH = "datasets/banking77_with_embeddings_test"
    MODEL_PATH = "random_forest_model.joblib"
    
    # Model parameters
    N_ESTIMATORS = 100
    MAX_DEPTH = None
    RANDOM_STATE = 42
    
    if TEST_MODE:
        print(f"⚠️  TEST MODE: Using only {TEST_SAMPLES} samples per split")
        print("Set TEST_MODE = False to use the full dataset")
    else:
        print("Using full dataset")
    
    try:
        # Step 1: Load dataset
        dataset = load_dataset_with_embeddings(DATASET_PATH)
        
        # Step 2: Prepare data
        X_train, y_train, X_test, y_test, label_names = prepare_data(
            dataset, test_mode=TEST_MODE, test_samples=TEST_SAMPLES
        )

        # Step 4: Train model
        model = train_random_forest(
            X_train, y_train, 
            n_estimators=N_ESTIMATORS,
            max_depth=MAX_DEPTH,
            random_state=RANDOM_STATE,
            test_mode=TEST_MODE
        )
        
        # Step 5: Evaluate model
        accuracy, predictions = evaluate_model(
            model, X_test, y_test, label_names, test_mode=TEST_MODE
        )
        
        # Step 6: Save model
        save_model(model, MODEL_PATH)
        
        # Step 7: Summary
        print("\n" + "="*60)
        print("SUMMARY")
        print("="*60)
        print(f"Model: Random Forest")
        print(f"Dataset: Banking77 with embeddings")
        print(f"Training samples: {len(X_train)}")
        print(f"Test samples: {len(X_test)}")
        print(f"Number of classes: {len(label_names)}")
        print(f"Embedding dimension: {X_train.shape[1]}")
        print(f"Test accuracy: {accuracy:.4f} ({accuracy*100:.2f}%)")
        print(f"Model saved to: {MODEL_PATH}")

        print("\n✅ Pipeline completed successfully!")
        
    except Exception as e:
        print(f"\n❌ Error occurred: {str(e)}")
        raise e

if __name__ == "__main__":
    main()

```
For this small example, we are not doing any hyperparameter tuning at all, which would improve the performance.

The model achieves 95% accuracy across 77 different banking service categories. The balanced nature of the dataset ensures that this 95% accuracy reflects genuine classification ability rather than simply predicting the most common class. With macro and weighted averages both at 0.95, the model demonstrates consistent performance across all banking intents, from card activations to account transfers, making it well-suited for real-world customer service applications.

These results are similar to the ones in [Fine-tune classifier with ModernBERT in 2025](https://www.philschmid.de/fine-tune-modern-bert-in-2025), showing the power of using a strong embedding model with a lightweight classifier.

