import torch
import os
import psutil
import time
import gc

from torch.optim import AdamW
from base_model.model import ActorModel, CriticModel
from utils.constant import DEVICE, SAVE_DIR

running_count = 0
running_mean = torch.zeros(7, dtype=torch.float, device=DEVICE)
running_std = torch.ones(7, dtype=torch.float, device=DEVICE)

def profile_with_psutil(func):
    """Profile using psutil for detailed system metrics"""
    def wrapper(*args, **kwargs):
        process = psutil.Process(os.getpid())

        # CPU baseline
        process.cpu_percent(interval=None)

        # Memory baseline
        mem_before = process.memory_info().rss

        # Run function
        result = func(*args, **kwargs)

        # Allow memory to stabilize
        time.sleep(0.05)
        gc.collect()

        # Memory after
        mem_after = process.memory_info().rss

        # CPU after
        cpu_after = process.cpu_percent(interval=0.1)

        mem_delta = mem_after - mem_before
        cpu_delta = cpu_after
        return result, mem_delta, cpu_delta
    return wrapper

def load_model(num_extra_feats):
    if not os.path.exists(SAVE_DIR + "/actor_model.pt") or not os.path.exists(SAVE_DIR + "/critic_model.pt"):
        actor = ActorModel(num_extra_feats = num_extra_feats)
        critic = CriticModel()
        return actor, critic
        
    actor_model = ActorModel(num_extra_feats = num_extra_feats)
    actor_model.load_state_dict(torch.load(SAVE_DIR + "/actor_model.pt", map_location=DEVICE))
    critic_model = CriticModel()
    critic_model.load_state_dict(torch.load(SAVE_DIR + "/critic_model.pt", map_location=DEVICE))

    return actor_model, critic_model

def extract_numeric_features(batch):
    feats = []
    keys = ["averageCyclomaticComplexity", "bundleSize", "totalCallCount", 
            "totalConditionalComplexity", "totalFunctionComplexity", 
            "totalLoopComplexity", "totalRecursionDepth"]
    for item in batch:
        feature_vec = []
        for i, k in enumerate(keys):
            val = item.get(k, 0.0)
            
            # Special handling for bundleSize - consider log transformation if it varies widely
            if k == "bundleSize" and val > 0:
                # Uncomment if bundleSize has wide range (e.g., 0.1MB to 100MB)
                # val = np.log10(val + 1)
                pass
                
            feature_vec.append(float(val))
        feats.append(feature_vec)

    features = torch.tensor(feats, dtype=torch.float, device=DEVICE)
    return normalize_features(features)


def normalize_features(features):
    """Normalize features using running mean and standard deviation."""
    global running_mean, running_std, running_count

    batch_size = features.shape[0]
    
    if running_count == 0:
        # Initialize with first batch
        running_mean = torch.mean(features, dim=0)
        running_std = torch.std(features, dim=0, unbiased=False)  # Use population std for consistency
        running_std = torch.clamp(running_std, min=1e-8)  # Better numerical stability
        running_count = batch_size
        
        normalized = (features - running_mean) / running_std
        return normalized

    # Welford's online algorithm for stable running statistics
    prev_count = running_count
    new_count = running_count + batch_size
    
    # Update running mean
    batch_mean = torch.mean(features, dim=0)
    delta = batch_mean - running_mean
    running_mean = running_mean + delta * batch_size / new_count
    
    # Update running std using Welford's method (more numerically stable)
    batch_var = torch.var(features, dim=0, unbiased=False)
    prev_var = running_std ** 2
    
    # Combine variances
    running_var = (prev_count * prev_var + batch_size * batch_var) / new_count + \
                  (prev_count * batch_size * delta ** 2) / (new_count ** 2)
    
    running_std = torch.sqrt(torch.clamp(running_var, min=1e-16))  # Prevent sqrt of negative
    running_count = new_count

    # Normalize current batch
    normalized = (features - running_mean) / running_std
    return normalized

def get_cpu_clock_rate_hz():
    """
    Retrieves the current CPU clock rate in Hertz.
    """
    try:
        # psutil.cpu_freq() returns a named tuple with current, min, and max frequencies in MHz.
        cpu_freq_info = psutil.cpu_freq()
        if cpu_freq_info:
            # Convert current frequency from MHz to Hz
            max_frequency_hz = cpu_freq_info.max * 1_000_000
            return max_frequency_hz
        else:
            return None
    except Exception as e:
        print(f"Error getting CPU frequency: {e}")
        return None