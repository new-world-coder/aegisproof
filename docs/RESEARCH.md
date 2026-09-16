# AegisProof — Phase 0 Research Notes

**Status:** Planning / research only. No claim of completed product functionality.  
**Date:** 2026-09-16

---

## 1. Research question

Is there a missing, vendor-neutral, developer-first interchange layer for **AI assurance evidence** — and if so, should it be a **new schema**, a **profile on existing standards**, or **tooling that speaks those standards**?

---

## 2. Gap analysis verdict (honest)

**The prompt’s premise that “there is no sufficiently simple evidence representation” is only partly true.**

What already exists:

| Layer | What exists | Maturity |
|---|---|---|
| AI/ML inventory (AIBOM / ML-BOM) | SPDX 3.0 AI + Dataset profiles; CycloneDX ML-BOM (ECMA-424) | Standardised, growing tooling |
| Claims + evidence + attestations | CycloneDX Declarations / Attestations (CDXA) | Spec-complete; adoption uneven |
| Security findings | SARIF | Widely adopted in CI |
| Supply-chain provenance | in-toto, SLSA, Sigstore, OCI referrers | Strong for build/deploy artefacts |
| Observability signals | OpenTelemetry GenAI, OpenInference, OpenLLMetry | Runtime telemetry, not assurance packs |
| Model/system documentation | Model cards, HF card metadata, Red Hat AI System Card schema | Fragmented, mostly descriptive |
| Assessment results (heavy GRC) | NIST OSCAL Assessment Results | Powerful, FedRAMP-heavy, poor DIY DX |
| EU AI Act local evidence packs | Regula (`regula-ai`) — hashed/signed packs mapped to articles | Direct CLI competitor for EU-specific path |
| Classification + assure + signed bundles | **RiskLine** (sister project) — already ships `assure` / `verify` + Ed25519 evidence | Overlaps hard with proposed AegisProof CLI |

**What is still weakly solved:**

1. A **developer-local, regulation-neutral evidence package** that *composes* BOM + eval artefacts + policies + CI outputs + optional runtime receipts into one verifyable unit.
2. A **SARIF-like gravity well** for AI assurance: one artefact CI emits and buyers/auditors/tools consume — without forcing OSCAL or full CycloneDX attestations on day one.
3. Clear separation of **evidence integrity** from **compliance determination** in a format engineers will actually author.

**Conclusion:** Building a greenfield “AI SBOM” or a GRC-style control checklist schema would **reinvent SPDX/CycloneDX and collide with Regula/RiskLine**. A thin **evidence package + verify CLI** that *references* those standards can still be justified — but only if it stays smaller than CycloneDX Attestations and does not become a second RiskLine.

---

## 3. Standards map (summary)

| Name | Solves | Does not solve for AegisProof | Reuse? |
|---|---|---|---|
| **SPDX 3.0 AI/Dataset** | Model/dataset inventory, licensing, relationships | Continuous control evidence packs; simple repo DX | **Reuse as inventory source** |
| **CycloneDX ML-BOM** | ML component inventory, model cards fields | Lightweight “verify this repo” DX | **Reuse as inventory / export target** |
| **CycloneDX Attestations** | Requirements → claims → evidence → conformance scores | Developer-first authoring UX; many teams find it heavy | **Primary interoperability target** |
| **SARIF** | Tool findings interchange | Assurance claims, model identity, control evidence | **Embed/reference for security evals** |
| **OpenAPI** | API contracts | Assurance evidence | Analogy only (DX + schema) |
| **OpenTelemetry / OpenInference / OpenLLMetry** | Runtime traces/metrics | Durable assurance artefact | **Reference runtime evidence URIs** |
| **in-toto / SLSA / Sigstore** | Build provenance, signing, transparency log | AI-specific claim taxonomy | **Later integrity upgrade path** |
| **W3C PROV** | Generic provenance graph | Practical CI packaging | Optional semantic mapping later |
| **OCI artifacts / referrers** | Attaching attestations to images/models | Repo-local evidence workflows | Optional distribution later |
| **MLflow / OpenLineage** | Experiment & job lineage | Assurance claim packaging | Reference as evidence locations |
| **Model / system cards** | Human documentation | Integrity, CI verify, claim↔artefact binding | Allow as evidence *types* |
| **OSCAL** | Control assessment results | Lightweight DX; AI-native claim types | Profile bridge later, not v0.1 core |
| **NIST AI RMF / ISO 42001 / 23894** | Process frameworks | Interchange schema | **Profiles**, not core schema |
| **EU AI Act** | Legal obligations, Annex IV docs | Neutral core schema | **Profile only** |
| **Regula evidence format** | EU AI Act hashed/signed packs | Regulation-neutral interchange | Compete carefully; interop later |
| **RiskLine evidence bundles** | Signed assure payloads, classify/assure/verify | Vendor-neutral multi-consumer format | **Extract/generalise or export to AegisProof/CDX** |

---

## 4. Reinvention risks (ordered)

1. **Duplicating CycloneDX Attestations** under a new brand.
2. **Duplicating RiskLine’s assure/verify/evidence** as a parallel product.
3. **Duplicating Regula** for EU AI Act evidence packs.
4. **Calling it “AI SBOM”** and fighting SPDX/OWASP for category ownership.
5. **Encoding EU AI Act into the core schema** (locks neutrality and adoption).
6. **Control-status fields (`implemented` / `verified`) without verification semantics** → compliance washing.
7. **Inventing another telemetry vocabulary** (lose to OTel GenAI / OpenInference).
8. **Custom signing/storage** instead of in-toto / Sigstore / OCI referrers later.

### Intellectual-honesty fork

If design converges on something isomorphic to **CycloneDX Declarations + ML-BOM subjects + in-toto signing**, the correct move is **contribute an AI assurance profile upstream** (or a normative mapping) rather than a parallel universe — unless a thinner CI-native package remains clearly cheaper for the wedge user.

Thread detail: [Standards research](3c740ebc-9407-4047-a64e-8a2c40e05959), [Competitive landscape](4424f005-cb0b-4f73-b771-cf392da808cd).

---

## 5. Smallest useful evidence object (research recommendation)

v0.1 should be closer to a **content-addressed evidence index** than a governance questionnaire:

```text
EvidencePackage
  identity: { system_id, version, producer }
  created_at
  items[]:
    id
    type          # evaluation | policy | sbom | sarif | test-result | model-card | runtime-receipt | other
    locator       # path or URI
    digest        # sha256:...
    media_type?   # optional
    produced_at?
    expires_at?   # freshness, not legal validity
  claims[]?       # optional in v0.1 — preferably deferred
    id
    statement     # free text or URI to requirement in a *profile*
    evidence_ids[]
```

**Non-goals for the smallest object:** risk classification, legal conformity, control matrices, vendor lock-in to RiskLine/XERJ, mandatory signing.

---

## 6. Adjacent OSS / commercial notes

- **Governance SaaS** (Credo, Holistic, OneTrust, watsonx.governance, Purview, ModelOp): dashboards and workflows — natural *consumers*, not peers.
- **Observability** (Langfuse, LangSmith, Phoenix): produce runtime signals AegisProof may *point to*, not replace.
- **Agent retrieval** (XERJ, Sourcegraph): index/search evidence; must not become the integrity root of trust.

---

## 7. Open research follow-ups

- [ ] Depth-read CycloneDX Attestations guide vs proposed package model (gap size quantification).
- [ ] Compare Regula `regula.evidence.v1` field-by-field with RiskLine `EvidenceBundle`.
- [ ] Survey how many AI teams already emit ML-BOM / SPDX AI in CI (adoption signal).
- [ ] Talk to 3 design partners: “Would you attach one file to a vendor questionnaire?”
- [ ] Decide: **new schema** vs **CycloneDX profile + DX sugar** (founding decision).

---

## 8. Sources (non-exhaustive)

- SPDX 3.0 AI/Dataset profiles; “Building an Open AIBOM Standard in the Wild” (experience report).
- CycloneDX ML-BOM guide; CycloneDX Attestations / ECMA-424.
- NIST OSCAL Assessment Results.
- Regula (`kuzivaai/getregula`) evidence pack + trust docs.
- RiskLine local repo: `pkg/evidence`, assure/verify CLI, ROADMAP.
- XERJ: xerj.org, reference-coding case study, ES-compatible autoindex.
- Red Hat AI System Card schema (early, low adoption signal).
