# Model Card (Synthetic)

## Overview
- **Model Name:** example-embeddings-v1
- **Provider:** Example AI Research (synthetic)
- **Model Type:** Dense bi-encoder text embedding model
- **Version:** 1.0.0
- **Context Length:** 512 tokens
- **Embedding Dimensions:** 768

## Intended Use
- **Primary Use Case:** Semantic retrieval over historical claims documentation and knowledge base articles for claims assistance (internal staff tool).
- **Out-of-Scope Use:** Direct automated decision-making, customer-facing chat generation, or medical diagnosis.

## Training & Evaluation Data
- **Training Data:** Synthetic corpus of insurance policy guidelines and mock claims narratives.
- **Evaluation Benchmark:** Synthetic claims retrieval test set (top-k=5 accuracy, faithfulness, recall).

## Factors & Limitations
- **Limitations:** Specialized for structured claim terminology; may degrade on unstructured conversational slang.
- **Known Biases:** Biased towards policy terminology present in synthetic benchmark distributions.

## Governance & Safety
- **PII Scrubbing:** All input documents pass through upstream redaction before embedding generation.
- **Integrity Tracking:** Managed via AegisProof manifest evidence bundle.
