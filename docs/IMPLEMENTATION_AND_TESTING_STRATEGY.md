# AegisProof — Implementation and Testing Strategy

**Phase 0 master planning document**  
**Status:** Founder approved (`approved` 2026-09-16) — Milestone 1–2 implementation underway / landed in-tree  
**Date:** 2026-09-16

Companion docs: `RESEARCH.md`, `ARCHITECTURE.md`, `SECURITY_STRATEGY.md`, `TESTING_AND_BENCHMARK.md`, `ROADMAP.md`, `OPEN_QUESTIONS.md`, `RISKLINE_INTEROPERABILITY.md`, `XERJ_INTEGRATION_HYPOTHESIS.md`, `WHITEPAPER_OUTLINE.md`.

---

## 1. Executive summary

AegisProof should explore whether **AI assurance evidence** can become a portable, vendor-neutral interchange primitive — closer to **SARIF + a package manifest** than to a GRC SaaS or an “AI SBOM.”

**Critical challenge to the founding prompt:** the gap is **real but narrower than claimed**. SPDX 3.0 AI profiles, CycloneDX ML-BOM + Attestations, Regula EU evidence packs, OSCAL, and **RiskLine’s existing assure/verify/signed evidence bundles** already cover large parts of the naive design. A new standard only earns the right to exist if it is:

1. **Thinner** than CycloneDX Attestations for day-1 developer DX,  
2. **Regulation-neutral** (unlike Regula’s EU-first packs),  
3. **Not a second RiskLine**, and  
4. Backed by a **default consumer** and measurable time savings.

**Recommended positioning:** *machine-readable evidence packages for AI systems* — not “AI SBOM.”

**Recommended v0.1 shape:** content-addressed evidence package (manifest + digests + freshness) + local `verify` CLI. No control-status theatre, no legal verdicts, no mandatory LLM/cloud/XERJ.

**Build status:** blocked on founder decisions in §26.

---

## 2. Problem definition

Organisations shipping AI systems must repeatedly assemble proof of *how* a system was built, tested, controlled, and operated. That proof today is scattered across:

source, config, model metadata, evaluations, policies, tests, CI, logs, monitoring, incidents, governance tools, deployment metadata, human approvals, security scanners.

Buyers and auditors ask for evidence; engineers produce screenshots and PDFs; platforms lock artefacts inside dashboards. The costly failure mode is **re-collection and re-formatting**, not lack of opinions about AI ethics.

---

## 3. Opportunity

If a **single, boring artefact** can be emitted in CI and consumed by multiple tools (vendor risk, RiskLine, scanners, agents, auditors), the industry gains an interoperability layer similar in *role* (not scale) to SARIF/SBOM files.

Economic wedge: reduce hours spent preparing assurance packets for enterprise procurement and internal audit — measurable, not rhetorical.

Social wedge: startups/SMEs/universities that cannot buy enterprise GRC still need a credible, inspectable evidence pack.

---

## 4. Strategic thesis

> AI systems need machine-readable, independently verifiable **evidence packages** describing artefacts of build, test, control, and operation — without conflating those packages with legal compliance.

Central strategic question:

> Can AI assurance evidence become an interoperability primitive rather than another proprietary dashboard?

**Thesis risk:** inventing a vanity standard that fragments an ecosystem already moving toward SPDX/CycloneDX.

**Thesis hedge:** ship DX + conformance on a thin package model with **first-class CycloneDX export**, and be willing to collapse into a CycloneDX profile project if that wins.

---

## 5. Existing ecosystem

See `docs/RESEARCH.md` for the full map.

**Reuse aggressively:**

- SPDX / CycloneDX for inventory  
- CycloneDX Attestations as primary heavy interop target  
- SARIF for security findings  
- SLSA / in-toto / Sigstore for later provenance  
- OTel/OpenInference as runtime evidence *pointers*  
- NIST/ISO/EU as **profiles**, not core schema  

**Do not reinvent:** AIBOM/ML-BOM fields, full attestation graphs, OSCAL programmes.

---

## 6. Competitor / adjacent analysis

### AI governance

| Actor | Relation |
|---|---|
| RiskLine | Sister engine (classify/assure); potential emitter/consumer — **must not own the standard** |
| Regula | OSS EU AI Act evidence CLI — closest **format+CLI** neighbour for EU path |
| Credo / Holistic / OneTrust / ModelOp / watsonx / Purview | Dashboard/workflow consumers |
| Fairlearn / RAI Toolbox | Evaluation tooling → evidence producers |

### AI observability

Langfuse, LangSmith, Phoenix, Helicone, OTel GenAI — produce signals; AegisProof may reference exports, not replace them.

### Supply chain

SPDX, CycloneDX, SLSA, Sigstore — **compose**.

### Agent retrieval

XERJ, Sourcegraph, vector DBs — optional indexers; **not trust roots**.

### Missing layer (precise)

Not “any AI metadata.” Missing is a **developer-first evidence package** that binds heterogeneous artefacts with integrity and travels across tools **without** requiring a SaaS or a single regulation ontology.

### Positioning candidates (ranked)

1. **Machine-readable evidence for AI systems** — clearest, least wrong  
2. **Open AI Evidence Standard** — good formal name once artefact exists  
3. **Open AI Assurance Evidence Layer** — sounds platform-y  
4. **AI SBOM** — **reject** as primary position (category owned; wrong mental model)

### Do not build if…

- Founder cannot name a non-RiskLine consumer who will reject releases without the file within ~2 milestones of CLI, **or**
- Decision is “full governance ontology in v0.1,” **or**
- Plan is to out-spec CycloneDX Attestations without tooling gravity.

---

## 7. RiskLine relationship

See `docs/RISKLINE_INTEROPERABILITY.md`.

RiskLine already has signed evidence verify. AegisProof should generalise the **portable package** problem; RiskLine keeps regime logic. Interop via import/emit only. Anti-coupling rules mandatory.

---

## 8. XERJ relationship

See `docs/XERJ_INTEGRATION_HYPOTHESIS.md`.

Exploratory only. Good for retrieval demos and reference-coding for *implementing* controls. Not strategically central to evidence integrity. Priority **P2** after spec+CLI.

---

## 9. Architecture

See `docs/ARCHITECTURE.md`.

Summary: local evidence package + verify CLI; compose BOM/SARIF/OTel; profiles for regimes; cloud optional; no `PARTIALLY ASSURED` compliance theatre in v0.1.

### Proposed repository structure (adjusted)

```text
aegisproof/
├── README.md
├── LICENSE
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── GOVERNANCE.md
├── SPEC.md                 # or /spec
├── schema/
├── cli/
├── examples/
├── conformance/
├── tests/
├── docs/                   # strategy + wiki sources
├── research/               # deep notes / bibliographies
├── integrations/           # riskline/, xerj/, github-action/ (later)
├── roadmap/
└── whitepaper/
```

Empty dirs from Phase 0 scaffolding may remain until Milestone 1.

---

## 10. Specification strategy

1. Write principles before fields.  
2. JSON Schema as source of truth; YAML as authoring sugar.  
3. v0.1 = package + evidence items + digests (+ optional freshness).  
4. Defer claims if they reopen compliance scope creep.  
5. Extension via `profiles[]` and open `type` strings with documented registry later.  
6. Versioning: `specVersion`; pre-1.0 breaking changes allowed with changelog.  
7. Conformance suite ships with schema — not afterthought.  
8. Parallel track: document CycloneDX Attestation mapping (even if export is Milestone 3+).

---

## 11. Security strategy

Full threat model: `docs/SECURITY_STRATEGY.md`. Summary below.

### Assets

- Integrity of evidence bytes  
- Trustworthiness of CI gate decisions  
- Developer trust that verify means what README says (and **not** more)

### Threats (priority)

| Threat | v0.1 mitigation |
|---|---|
| Path traversal via locators | Resolve under package root; reject `..` and absolute paths |
| Malicious YAML / billion laughs | Prefer JSON for machine path; YAML via safe loader; size caps |
| Oversized inputs | Max bytes / max files |
| Hash mismatch / substitution | Required digests; fail closed |
| Stale evidence | Optional `expiresAt`; explicit exit/warn policy |
| Dependency supply chain on CLI | Minimal deps; pin; later provenance for releases |
| Compliance washing | UX copy: never “compliant”; README disclaimer tested |
| Signature theatre | No signing until verify semantics are boring and correct |
| Untrusted remote URIs | Local-only locators in v0.1 |

### Honest residual risk

A green verify means **package integrity**, not model safety, not legal conformity, not absence of fraud by the producer.

### Disclosure

`SECURITY.md` with private reporting channel; 90-day disclosure aim; no bounty required at start.

### Later upgrades

Ed25519 or Sigstore; RFC3161 timestamps; in-toto placements — after hash path is solid.

---

## 12. Testing strategy

Detailed methodology (Northwind pilot, kill criteria, vanity policy): `docs/TESTING_AND_BENCHMARK.md`.

### 12.1 Schema / conformance

- Valid packages  
- Missing required fields  
- Unknown fields policy (forbid vs allow)  
- Extension types  
- Back-compat fixtures per `specVersion`

### 12.2 Evidence integrity

- Missing file  
- Bad digest  
- Broken reference  
- Modified after pack  
- Stale `expiresAt`  
- Empty evidence list policy

### 12.3 Security tests

- `../` escape attempts  
- Symlink escape  
- Huge file / huge manifest  
- Hostile YAML constructs  
- Nested archives if ever supported (initially: don’t)

### 12.4 Developer experience metrics

Track in `conformance/dx.md` over time:

| Metric | Target (proposal) |
|---|---|
| Install time | &lt; 2 minutes |
| Time to first verify on example | &lt; 10 minutes |
| Files required for minimal package | ≤ 2 (manifest + 1 evidence) |
| CI integration | ≤ 20 lines workflow |

### 12.5 Interoperability tests

- GitHub Actions  
- JSON and YAML manifests  
- Fixture export mapping to CycloneDX (when implemented)  
- Contract tests for future RiskLine import (golden files)  
- XERJ indexing smoke (optional job, non-blocking)

---

## 13. Benchmark methodology (empirical impact)

Canonical protocol: **`docs/TESTING_AND_BENCHMARK.md`** (Northwind Claims Assistant, 15-item checklist, A→B hard caps, M1–M6 metrics).

**Success (continue investing):** meaningful reduction in `(engineer + review) hours` **or** fewer missing checklist items with integrity verify available.

**Hard kill:** After not better on hours *and* missing items; or integrity false-negatives; or TTFV &gt; 45 min after one doc pass; or CDX/SPDX already cover the checklist cheaper → pivot to profile.

Publish method and raw timings; allow null results. n=1 = pilot only; ≥3 pilots before strong public impact claims.

---

## 14. Cloud strategy

**Principle:** local-first, cloud-optional.

Phase 0–2: no cloud services required.

Later reference architecture (only if demand):

```text
CLI/SDK → optional service → evidence store / policy engine / audit API
```

Prefer: object storage + serverless verify + static docs. Avoid always-on Kubernetes for a single maintainer. Docs: GitHub Pages is enough initially.

---

## 15. Repository structure

See §9. Licence recommendation: **Apache-2.0** (aligns with RiskLine, XERJ, broad corporate acceptance). Confirm trademark/name search before wide marketing.

---

## 16. Milestones

See `docs/ROADMAP.md` (M0–M8). First three after approval:

1. Minimal specification + conformance  
2. CLI verify  
3. Integrity hardening  

---

## 17. Open questions

See `docs/OPEN_QUESTIONS.md`.

---

## 18. Risks

| Risk | Severity | Mitigation |
|---|---|---|
| Reinvention of CycloneDX/SPDX | High | Export/mapping; pivot-to-profile option |
| Collision with RiskLine evidence | High | Explicit split; anti-coupling |
| Collision with Regula (EU) | Medium | Neutral core; interop later |
| No default consumer | High | Wedge = vendor-risk packet; seek design partner |
| Compliance washing UX | High | Forbidden verdict language; tests on disclaimer |
| Over-scope v0.1 | High | Evidence items only |
| XERJ distraction | Medium | P2 experiments only |
| Vanity metrics | Medium | Benchmark kill criteria |
| Solo-maintainer burnout | Medium | Milestone caps; no cloud early |
| Adoption failure | High | 12–18 month kill/fold criteria |

---

## 19. Open-source governance

**Early (now → multi-emitter):**

- Benevolent maintainer + public issues  
- DCO (Developer Certificate of Origin) on commits  
- Apache-2.0  
- `GOVERNANCE.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`  
- Spec changes via proposal issues / lightweight RFC  

**Later (2+ independent emitters):**

- Technical steering with emitter representatives  
- Compatibility policy for ≥1.0  

**Foundation:** only if neutrality becomes a real adoption blocker — not as cosplay.

**Vendor neutrality:** RiskLine is one implementation; cannot hold exclusive advance notice on breaking changes beyond normal maintainer coordination.

---

## 20. Whitepaper plan

See `docs/WHITEPAPER_OUTLINE.md`. Draft after Milestone 1 freezes vocabulary. No partnership claims. Include “Why not only CycloneDX?”

---

## 21. Wiki plan

Planned sections:

1. What is AegisProof?  
2. Why it exists (and what it is not)  
3. Architecture  
4. Specification  
5. Evidence model  
6. Security  
7. Integrations (CI, RiskLine, XERJ experiments)  
8. Regulatory profiles  
9. FAQ (esp. evidence ≠ compliance)  
10. Governance  
11. Contributing  

README must explain the problem in &lt;2 minutes; wiki holds depth.

---

## 22. GitHub issue backlog (initial)

Labels: `spec` `architecture` `research` `security` `cli` `testing` `documentation` `integrations` `riskline` `xerj` `good-first-issue` `discussion`

### Spec

- Define evidence object model (minimal)  
- Choose schema format (JSON Schema) + YAML authoring rules  
- Evidence type registry (v0 draft)  
- Extension / profile mechanism  
- Versioning and compatibility policy  
- Model identity: reference ML-BOM vs duplicate fields  
- Evaluation representation (pointer vs embedded)  
- Freshness / expiry semantics  

### Architecture

- Dual YAML vs CycloneDX-native decision record  
- Package directory convention  
- Exit code policy  

### Research

- Field-level mapping: CycloneDX Attestations ↔ AegisProof  
- Field-level mapping: Regula evidence v1 ↔ AegisProof  
- Field-level mapping: RiskLine EvidenceBundle ↔ AegisProof  
- Survey: do design partners already generate ML-BOM?  

### Security

- Threat model acceptance  
- Path confinement rules  
- Input size limits  
- Disclosure process  

### CLI

- `verify` UX + JSON schema for report  
- `init` / digest fill  
- Language/runtime choice for CLI  

### Testing

- Conformance suite skeleton  
- Security fixture pack  
- Empirical benchmark fixture “claims-assistant”  

### Documentation

- README draft (&lt;2 min)  
- Disclaimer language tests  
- Wiki stubs  

### Integrations

- GitHub Action design  
- CycloneDX export spike  

### RiskLine

- Import experiment design  
- Emit experiment design  
- Anti-coupling checklist in CI of *this* repo  

### XERJ

- Experiment 1–4 tracking issues  

### Good first issues

- Add invalid conformance fixture  
- Improve example model card evidence  
- Docs typo / glossary  

### Discussion

- Positioning: reject AI SBOM?  
- Should claims exist in v0.1?  
- Kill criteria acceptance  

---

## 23. Future roadmap

See `docs/ROADMAP.md`. Profiles, signing, registries, federation, agent interfaces — **conditional** on adoption signals.

---

## 24. Success metrics

### Adoption (secondary early, primary later)

Stars/forks/downloads are **vanity until** emitters exist. Track but do not optimise.

### Leading indicators (what matters)

- # independent emitters (non-founder, non-RiskLine)  
- # consumers that gate on verify  
- Time-to-first-verify  
- Benchmark time savings  
- Conformance suite external runs  

### Integration

- RiskLine import/emit  
- CI action usage  
- Optional XERJ experiment completion  

### Economic / social

- Hours saved in benchmark + design-partner anecdotes  
- Use by SME/university fixtures  

### Kill criteria

Within 12–18 months of first CLI release: no third-party emitter **or** no recurring consumer → fold into RiskLine adapters / CycloneDX profile effort.

---

## 25. Career / technical portfolio objectives

GitHub history should read:

1. Elasticsearch JS — production OSS infrastructure contribution  
2. RiskLine — deterministic AI risk/assurance engineering  
3. AegisProof — **neutral interchange** / standards-shaped systems thinking  

Demonstrate: schema design, security threat modelling, CI tooling, provenance, DX, governance — **not** vanity dashboards.

Avoid: pairing AegisProof releases exclusively to RiskLine marketing screenshots.

---

## 26. Explicit decisions requiring founder approval

| ID | Decision | Recommendation |
|---|---|---|
| D1 | Proceed with AegisProof as separate project? | **Yes, but thin** — or pivot to “CycloneDX AI evidence DX” if D2=B |
| D2 | Schema strategy: A dual YAML+CDX export / B CDX-native / C independent-only | **A** |
| D3 | Reject “AI SBOM” as primary positioning? | **Yes** |
| D4 | v0.1 exclude control statuses & risk classification? | **Yes** |
| D5 | RiskLine: extract format vs independent-then-adapt? | **Independent-then-adapt**, with RiskLine as first *optional* emitter — plus **one non-RiskLine example required** |
| D6 | CLI language | **Go** (ops story aligned with RiskLine) or Python (faster DX) — founder pick |
| D7 | Licence | **Apache-2.0** |
| D8 | Claims in v0.1? | **Defer** |
| D9 | Hosted cloud before traction? | **No** |
| D10 | XERJ in public architecture diagrams? | **No** until experiments done |
| D11 | Accept kill/fold criteria §24? | **Yes** |
| D12 | Project org: personal vs new org | New `aegisproof` org preferred for neutrality optics |
| D13 | Name “AegisProof” trademark/collision check | Run before launch noise |
| D14 | Relationship messaging with Ivan/XERJ | Technical experiment language only |

---

## Business model (timing note)

Free forever: spec, schemas, CLI, basic verify, examples, conformance.

Commercial *later* (only after pull demand): hosted retention, enterprise support, conformance help, integration engineering. **Not** monetising the schema itself. RiskLine remains the natural commercial control plane.

---

## FOUNDER APPROVAL REQUEST

### Proposed project

**AegisProof** — open, vendor-neutral evidence package format + reference CLI for AI systems.

### One-sentence purpose

Make AI assurance evidence portable and verifiable in CI without turning it into a compliance verdict or a SaaS dashboard.

### Core user

Platform / AppSec / ML engineer who must attach credible artefacts to enterprise security / vendor-risk processes on every release.

### Core problem

Evidence is fragmented and non-interoperable; re-assembly is expensive; dashboards do not travel.

### Proposed v0.1

JSON Schema for an evidence package (identity + evidence items with locators/digests/freshness) + examples + conformance fixtures + local `aegisproof verify` (human/JSON/exit codes). No regime engine, no control-status matrix, no cloud, no LLM.

### Architecture

Local-first package → verify CLI → consumers (CI, auditors, optional RiskLine, optional indexers). Compose SPDX/CycloneDX/SARIF; profiles later for EU/NIST/etc.

### What is deliberately excluded

GRC SaaS, legal compliance determination, AI SBOM reinvention, mandatory XERJ, mandatory RiskLine, control “implemented/verified” theatre, early crypto theatre, early cloud.

### Key standards reused

SPDX, CycloneDX (ML-BOM + Attestations mapping), SARIF, later SLSA/Sigstore; OTel as pointers; regimes as profiles.

### RiskLine relationship

Sister consumer/emitter; not owner of the standard; no core dependency.

### XERJ relationship

Optional retrieval experiment only; not a partner claim; not a trust root.

### Proposed repository structure

As in §9 / `docs/ARCHITECTURE.md`.

### First 3 milestones

1. Minimal spec + conformance  
2. CLI verify  
3. Integrity hardening  

### Biggest technical risks

Reinvention vs CycloneDX; overlap with RiskLine evidence; unsafe parsers/path traversal; compliance-washing UX.

### Biggest adoption risks

No default consumer; only RiskLine uses it; “another YAML”; EU-specific tools (Regula) occupying mindshare.

### Estimated implementation complexity

**M1–M2:** small (days–low weeks for experienced founder) if scope stays thin.  
**M3–M5:** moderate.  
**Standards politics / interop:** harder than code.

### Questions requiring founder decisions

**D1–D14** in §26.

---

> **BUILD STATUS: UNBLOCKED — MILESTONE 1–2 IN TREE**

Defaults applied from §26 recommendations (D2=A dual YAML, D3 reject AI SBOM, D4 no control statuses, D6 Go CLI, D7 Apache-2.0, D8 defer claims, D9 no cloud, D10 XERJ off critical path).
