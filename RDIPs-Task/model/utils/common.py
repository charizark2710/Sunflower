import torch
from typing import Dict, List

DEVICE = 'cuda' if torch.cuda.is_available() else 'cpu'
MAX_TOKEN_LEN = 512
CPU_RATE = 3.5e9  # Example CPU rate in cycles per second

def tokenize_codes(tokenizer ,codes: List[str]) -> tuple[List[torch.Tensor], List[torch.Tensor]]:
    input_ids_list = []
    attn_masks_list = []

    for code in codes:
        # split code string by semicolon but keep the semicolon at the end
        parts = [p.strip() + ";" for p in code.split(";") if p.strip()]

        chunks = []
        cur_chunk = ""

        final = False
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
                final = True
            else:
                # finalize the current chunk
                if cur_chunk:
                    final = False
                    chunks.append(cur_chunk)
                else: final = True
                cur_chunk = part

        if final == True:
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
