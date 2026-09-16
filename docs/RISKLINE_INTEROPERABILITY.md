# AegisProof ↔ RiskLine Interoperability (Phase 0)

**Status:** Planning hypothesis. No integration code.  
**Constraint:** AegisProof must remain vendor-neutral; RiskLine must not own the standard.

---

## 1. Current RiskLine reality (important)

As of local `riskline` tree review (2026-09):

RiskLine already ships:

- Deterministic **classify** (EU AI Act / NIST packs)
- **Assure** layer with SHA-256 evidence records
- CLI **`verify`** for Ed25519-signed evidence bundles (`pkg/evidence`)
- CI GitHub Action path (`riskline-assure`)
- Runtime register/verify direction on the roadmap

**Implication:** A large fraction of the prompt’s “Milestone 2–3 CLI” already exists *inside RiskLine*. Creating a parallel `aegisproof verify` that only understands RiskLine payloads would be product theatre.

---

## 2. Desired split of concerns

| Concern | Owner |
|---|---|
| Regime rules, classification, control evaluation | **RiskLine** |
| Neutral evidence package interchange + integrity verify | **AegisProof** |
| EU-specific article packs (adjacent OSS) | Regula / profiles |
| Inventory BOM | SPDX / CycloneDX |
| Hosted dashboards / enterprise SaaS | RiskLine commercial path (optional) |

```text
AegisProof  →  open evidence format
     ↓
RiskLine    →  consumes + optionally emits
     ↓
Governance controls / continuous verification / dashboards
```

---

## 3. Interop model (proposed)

### 3.1 RiskLine as consumer

```bash
# future — not implemented
riskline import-aegisproof ./aegisproof.yaml
riskline assure --evidence-package ./aegisproof.yaml
```

RiskLine maps package evidence items into assure inputs / probe results where types align; unknown types are preserved as opaque attachments.

### 3.2 RiskLine as emitter

```bash
# future — not implemented
riskline assure ... --emit-aegisproof ./out/
```

Emits a regulation-neutral package: digests of classification JSON, assure response, probe files — **without** claiming AegisProof “validates compliance.”

### 3.3 Format boundary

- RiskLine `EvidenceBundle` (signed AssureResponse) remains a **RiskLine artefact**.
- AegisProof package may *contain* that bundle as an evidence item (`type: riskline-assure-bundle` or generic `application/vnd.riskline.evidence+json`).
- Core AegisProof schema must not require RiskLine types.

---

## 4. Anti-coupling rules

1. AegisProof repo must not import RiskLine as a dependency for core verify.
2. README must not say “powered by RiskLine.”
3. Spec changes require AegisProof maintainers; RiskLine adapts.
4. Trademark / branding: “RiskLine-compatible” allowed; “official RiskLine standard” forbidden.
5. If only RiskLine ever emits/consumes AegisProof for 12–18 months → **kill or fold** (see success metrics).

---

## 5. What belongs in RiskLine, not AegisProof

- EU / NIST / MAS regime logic
- Conformity state machines (`green` / material change)
- Runtime observation policy
- GRC dashboards
- “Recommended controls” narrative

---

## 6. What belongs in AegisProof, not RiskLine

- Portable package schema + conformance tests
- Path-safe digest verification
- Multi-emitter examples (plain RAG app, agent, eval CI) with **zero** RiskLine dependency
- Profile mechanism for regulatory mappings as *data*, not engine

---

## 7. Migration / extraction option

If founder chooses maximum honesty:

> Extract a neutral subset of RiskLine’s evidence hashing/signing into AegisProof, then make RiskLine depend on the open format.

This is strategically clean but slows RiskLine feature velocity. Requires explicit approval.

---

## 8. Experiment plan (post-approval, still not Phase 0 code)

1. Manually author an AegisProof package for a RiskLine example system.
2. Document field mapping table: AssureResponse ↔ evidence items.
3. Prototype import in a **separate** `integrations/riskline` folder or sibling repo.
4. Measure: does import reduce questionnaire prep time vs raw RiskLine JSON?

---

## 9. Founder decision needed

- **I1:** Extract evidence format from RiskLine vs define AegisProof independently then adapt?
- **I2:** Is RiskLine allowed to be the *first* emitter, or must a non-RiskLine example ship on day one?
- **I3:** Should AegisProof absorb RiskLine’s Ed25519 bundle format or treat it as opaque evidence?
