# AegisProof Whitepaper Outline (Phase 0)

**Status:** Outline only. Do not publish claims of adoption, partnerships, or legal effect.

Working title options:

1. *AegisProof: Machine-Readable Evidence Packages for AI Systems*
2. *Toward an Interoperability Layer for AI Assurance Evidence* (safer, less brand-forward)

---

## 1. Abstract

- Problem: AI assurance evidence is fragmented across repos, CI, evals, policies, and dashboards.
- Hypothesis: a thin, vendor-neutral evidence **package** (not a full GRC system) can become an interchange primitive.
- Method: compose existing standards; measure developer time-to-assemble evidence.
- Non-claim: not “the OpenTelemetry of AI,” not legal compliance.

## 2. Problem

- Evidence archaeology for vendor risk / audits.
- Non-comparable artefacts across teams.
- Dashboards that do not travel.
- Confusion between inventory (BOM), telemetry (OTel), findings (SARIF), and assurance claims.

## 3. Historical open infrastructure analogy (careful)

- What TCP/IP, Git, OpenSSL, Let’s Encrypt, SPDX, SARIF actually standardised.
- Strongest analogy: **SARIF** (tool evidence interchange) + **SBOM demand pattern**.
- Weakest analogy: Kubernetes/CNCF day-zero; OpenTelemetry-as-platform.

## 4. AI assurance problem

- Systems, not only models.
- Continuous change (prompts, tools, RAG corpora, models).
- Evidence ≠ compliance.
- Runtime drift vs design-time docs.

## 5. Existing ecosystem

- SPDX / CycloneDX ML-BOM & Attestations
- SARIF, SLSA, in-toto, Sigstore, OCI
- OTel GenAI / OpenInference
- Model cards / system cards / OSCAL
- Regula EU evidence packs
- RiskLine classify/assure/verify
- Governance SaaS as consumers

**Honest gap statement:** inventory and heavy attestation exist; lightweight package DX and multi-emitter gravity are weak.

## 6. Design principles

- Open first, vendor neutral, no mandatory LLM
- Machine-readable + human-inspectable
- Cryptographically useful without crypto theatre
- Composable; regulation-neutral core; profiles for regimes

## 7. Evidence model

- Package + evidence items + digests + freshness
- Optional claims
- Explicit exclusion of control-status theatre in core

## 8. Architecture

- Local-first CLI
- Consumers: CI, RiskLine, auditors, optional indexers
- Cloud-optional later

## 9. Security

- Threat model summary
- Path traversal, malicious YAML, hash substitution, compliance washing
- Staged integrity: hash → sign → Sigstore

## 10. Provenance

- Content digests as root
- Git refs
- Relationship to SLSA/in-toto

## 11. Interoperability

- Import/export with CycloneDX Attestations
- SARIF/SPDX as evidence types
- Profile mappings

## 12. RiskLine relationship

- Sister consumer/emitter
- Non-ownership of standard
- What stays in RiskLine

## 13. XERJ relationship

- Optional retrieval experiment
- Not root of trust
- Reference-coding vs evidence verify

## 14. Evaluation methodology

- Conformance suites
- Security tests
- DX metrics

## 15. Economic impact methodology

- Before/after engineer & governance hours
- Missing evidence counts
- Duplicate format counts
- Kill criteria if no savings

## 16. Limitations

- Will not replace lawyers, notified bodies, or OSCAL programmes
- Will not invent trust where producers lie
- Green verify ≠ safe system

## 17. Governance

- Maintainer model → multi-emitter TSC → possible foundation later
- Compatibility guarantees at 1.0

## 18. Future roadmap

- Milestone summary from `docs/ROADMAP.md`
- Profiles, signing, registries — conditional on adoption

---

## Writing rules

- Cite primary sources; no partnership claims.
- Prefer measured language (“we propose,” “we hypothesise”).
- Include a “Why not only CycloneDX?” section with a reversible decision.
- Target length for first draft: 8–12 pages, not a book.
