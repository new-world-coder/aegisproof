# AegisProof — Roadmap (Phase 0 Proposal)

**Status:** Milestone-based. No calendar dates. Subject to founder approval.  
**Rule:** Do not start Milestone 1 until explicit `APPROVED — START BUILD`.

---

## Guiding constraints

- Local-first, no mandatory LLM, no mandatory cloud.
- Evidence ≠ compliance.
- Regulation-neutral core; profiles later.
- Prefer reuse of SPDX / CycloneDX / SARIF / Sigstore over reinvention.
- XERJ and RiskLine integrations are experiments, not identity.

---

## Milestone 0 — Research & architecture ✅

Deliverables completed in `docs/`.

## Milestone 1 — Minimal specification ✅ (v0.1 draft)

- `SPEC.md`, `schema/aegisproof-0.1.schema.json`
- `examples/minimal`, `examples/rag-claims`
- `conformance/v0.1`

## Milestone 2 — CLI verifier ✅ (initial)

- `cmd/aegisproof`: `verify`, `pack`, `init`
- Exit codes 0/1/2/3; human + JSON
- Path confinement + SHA-256 digests
- Unit tests under `internal/`

## Milestone 3 — Integrity hardening

- Canonicalisation notes
- Optional Ed25519 signing **or** documented path to Sigstore
- Git commit binding optional
- Threat-model mitigations from Phase 0 security section

## Milestone 4 — CI

- GitHub Action
- Generic CI examples (GitLab/Jenkins snippets)
- Fail on `INVALID` default

## Milestone 5 — Reference implementations

- Simple AI app
- RAG app (example exists; expand)
- Agent
- “High-risk-style” fixture (synthetic — no legal claim)

Each emits package in CI.

## Milestone 6 — RiskLine interop experiment

- Separate integration doc/PR
- Import and/or emit
- No core dependency

## Milestone 7 — XERJ experiment

- Run four experiments in `docs/XERJ_INTEGRATION_HYPOTHESIS.md`
- Publish null results if null

## Milestone 8 — Ecosystem SDKs

- Python / TypeScript / Java **only after** schema stability signal
- Prefer generating from JSON Schema

---

## Version labels (spec)

| Label | Meaning |
|---|---|
| **v0.1** | Minimal package schema + verify _(current draft)_ |
| **v0.2** | Integrity / optional signing |
| **v0.3** | CI integrations |
| **v0.4** | Reference implementations |
| **v0.5** | Interop experiments |
| **v1.0** | Stable spec + compatibility guarantees |

---

## Explicitly deferred / maybe never

- Hosted evidence registry
- Compliance certification product
- Foundation paperwork before multi-emitter reality
- Regulatory profiles as paid-only artefacts (profiles should stay open; *support* may be commercial later)

---

## Kill / pivot triggers

- CycloneDX Attestations + thin CLI wrapper satisfies wedge users → **pivot to CDX profile project**.
- Only RiskLine uses AegisProof after 12–18 months → **fold into RiskLine adapters**.
- Benchmarks show no audit-prep time savings → **stop expanding schema**; reassess wedge.
- Name/positioning collides badly → rename rather than force “AI SBOM.”
