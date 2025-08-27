import torch
from typing import Dict, List


DEVICE = "cpu"
MAX_TOKEN_LEN = 512

num_features = 7
running_mean = torch.zeros(num_features, dtype=torch.float, device=DEVICE)
running_std = torch.ones(num_features, dtype=torch.float, device=DEVICE)
running_count = 0


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

# def tokenize_codes(codes: List[str]): 
#     global tokenizer 
#     enc = tokenizer(codes, padding=True, truncation=True, max_length=MAX_TOKEN_LEN, return_tensors="pt") 
#     return enc["input_ids"], enc["attention_mask"]

def tokenize_codes(tokenizer ,codes: List[str]) -> tuple[List[torch.Tensor], List[torch.Tensor]]:
    input_ids_list = []
    attn_masks_list = []

    for code in codes:
        # split code string by semicolon but keep the semicolon at the end
        parts = [p.strip() + ";" for p in code.split(";") if p.strip()]

        chunks = []
        cur_chunk = ""

        for part in parts:
            # tentative chunk if we add this part
            tentative = cur_chunk + " " + part if cur_chunk else part

            # check token length
            tokenized = tokenizer.encode(
                tentative,
                add_special_tokens=False
            )

            if len(tokenized) <= (MAX_TOKEN_LEN - 2):
                cur_chunk = tentative
            else:
                # finalize the current chunk
                if cur_chunk:
                    chunks.append(cur_chunk)
                # start new chunk with current part
                cur_chunk = part

        if cur_chunk:
            chunks.append(cur_chunk)

        # Now tokenize each chunk and pad
        chunk_tensors = []
        mask_tensors = []

        for ch in chunks:
            ids = tokenizer.encode(ch, add_special_tokens=False)[: MAX_TOKEN_LEN - 2]
            cur = [tokenizer.cls_token_id] + ids + [tokenizer.sep_token_id]
            attn = [1] * len(cur)

            pad_len = MAX_TOKEN_LEN - len(cur)
            if pad_len > 0:
                cur += [tokenizer.pad_token_id] * pad_len
                attn += [0] * pad_len

            chunk_tensors.append(torch.tensor(cur, dtype=torch.long))
            mask_tensors.append(torch.tensor(attn, dtype=torch.long))

        input_ids_list.append(torch.stack(chunk_tensors, dim=0))
        attn_masks_list.append(torch.stack(mask_tensors, dim=0))

    return input_ids_list, attn_masks_list

