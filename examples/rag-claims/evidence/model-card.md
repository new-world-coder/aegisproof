# Model card (synthetic example)

> This file is fictional and exists only to demonstrate an AegisProof evidence
> package. It must not be used to describe a deployed model or to support a
> safety, compliance, or procurement decision.

## Model details

- Provider: Example AI Research
- Identifier: example-embeddings-v1
- Version: 1.4.2-example
- Modality: text embeddings
- Embedding dimensions: 768
- Input limit: 512 tokens
- Release status: illustrative only; not a real model release

## Intended use

The example system uses this model to retrieve passages for a staff-facing
claims-assistance workflow. Retrieval results are supporting context only;
staff remain responsible for reviewing source material and making decisions.

## Out of scope

- Automated approval, denial, or ranking of claims
- Biometric identification or inference
- Use with personal or production data

## Example limitations

- The example does not establish retrieval quality, fairness, privacy, or
  safety for any real-world deployment.
- Short, ambiguous, or domain-specific queries can retrieve irrelevant
  passages.
- Retrieved text can be incomplete, outdated, or unsupported by the source
  record; it should not be treated as a factual answer.

## Example monitoring and review

For a real system, owners would record the embedding model version, retrieval
corpus version, evaluation results, and material incidents for each release.
Those records are intentionally absent from this synthetic example.
