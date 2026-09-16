# AegisProof

**Machine-readable evidence packages for AI systems** — local verify, no LLM required.

> Integrity check only. `VALID` does **not** mean compliant, certified, assured, or safe.

## Quick start

```bash
go run ./cmd/aegisproof verify ./examples/minimal
go run ./cmd/aegisproof verify ./examples/minimal --format json
```

Scaffold a package:

```bash
go run ./cmd/aegisproof init ./my-system
# add evidence files…
go run ./cmd/aegisproof pack ./my-system
go run ./cmd/aegisproof verify ./my-system
```

## What this is

A thin **evidence package** format (`aegisproof.yaml`) that lists artefacts (evals, policies, tests, SBOMs, …) with SHA-256 digests, plus a CLI that checks they match.

Think **SARIF-shaped interchange for assurance artefacts**, not an “AI SBOM” and not a GRC product.

## What this is not

- Legal compliance determination  
- Risk classification (see RiskLine as a possible consumer)  
- A dashboard or cloud service  
- A mandatory dependency on any vendor or indexer  

## Spec & schema

- [`SPEC.md`](SPEC.md) — v0.1 specification  
- [`schema/aegisproof-0.1.schema.json`](schema/aegisproof-0.1.schema.json)  
- [`conformance/`](conformance/) — valid/invalid fixtures  
- [`examples/`](examples/) — minimal + RAG samples  

## Exit codes

| Code | Result |
|---|---|
| 0 | `VALID` |
| 1 | `INVALID` |
| 2 | `STALE` (digests OK; freshness warnings) |
| 3 | Usage / IO error |

## Phase 0 planning docs

See [`docs/IMPLEMENTATION_AND_TESTING_STRATEGY.md`](docs/IMPLEMENTATION_AND_TESTING_STRATEGY.md) and siblings under `docs/`.

## Licence

Apache-2.0 — see [`LICENSE`](LICENSE).
