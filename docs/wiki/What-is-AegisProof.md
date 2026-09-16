# What is AegisProof?

AegisProof is an **open, vendor-neutral evidence package format** plus a local CLI.

You place an `aegisproof.yaml` (or `.json`) in a project. It lists artefacts about an AI system revision — evaluations, policies, tests, SBOMs, model cards, runtime receipts — each with a SHA-256 digest. `aegisproof verify` checks that those files exist and match their digests.

It is **not** a compliance verdict, risk score, or GRC dashboard.
