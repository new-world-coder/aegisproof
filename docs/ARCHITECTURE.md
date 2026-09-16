# AegisProof — Conceptual Architecture (Phase 0)

**Status:** Design proposal for founder approval. Not implemented.

---

## 1. One-line architecture thesis

AegisProof should be a **local-first evidence package format + verifier**, not an inventory standard, not a compliance engine, and not a governance platform.

```text
AI SYSTEM ARTEFACTS (repos, CI, evals, policies, BOMs, runtime receipts)
        │
        ▼
  Evidence Package Manifest  (aegisproof.yaml / .json)
        │  references + digests + freshness
        ▼
  aegisproof verify   →  human report + JSON + exit codes
        │
        ├── CI / AppSec
        ├── RiskLine (consumer / optional emitter)
        ├── Auditors / vendor risk
        └── Optional indexers (e.g. XERJ) — never root of trust
```

---

## 2. Challenge to the illustrative YAML

The founder prompt’s sample mixes three layers that should stay separate:

| Layer in sample | Problem |
|---|---|
| `system` / `model` inventory | Partially duplicates SPDX/CycloneDX ML-BOM |
| `controls.*.status: implemented/verified` | **Compliance washing** — status without verification semantics |
| `risk.classification: high` | Belongs in RiskLine / regime engines, not the evidence core |
| `evidence[]` locs | This is the only part that should be core |

**v0.1 recommendation:** drop control statuses and risk classification from the core schema. Profiles may *map* evidence to controls later; the core only proves that referenced artefacts exist, match digests, and meet freshness rules.

---

## 3. Evidence object model (proposed v0.1)

### 3.1 Package (manifest)

Required:

- `specVersion`
- `system` — `{ name, version, id? }`
- `producer` — tool/org that created the package
- `createdAt`
- `evidence[]` — list of evidence items

Optional in v0.1:

- `claims[]` — deferred preferred; if included, claims are **assertions**, not verdicts
- `profiles[]` — e.g. `eu-ai-act@draft` as opaque tags (no legal meaning)
- `externalDocuments[]` — pointers to SPDX/CycloneDX BOMs

### 3.2 Evidence item (smallest useful object)

```yaml
id: safety-eval-2026-09
type: evaluation          # enum + extension string
locator: ./evidence/safety.json
digest: sha256:…
mediaType: application/json
producedAt: 2026-09-01T12:00:00Z
expiresAt: 2026-12-01T00:00:00Z   # optional freshness
```

**Verify semantics (v0.1):**

1. Schema-valid manifest.
2. Every locator resolves (local path within package root; no path escape).
3. Digest matches file bytes.
4. Optional: warn/fail if `expiresAt` passed.
5. Never assert “compliant” or “assured” as a legal state.

Suggested result vocabulary (deliberately boring):

- `VALID` — package integrity OK
- `INVALID` — schema/hash/path failure
- `STALE` — integrity OK but freshness warnings (exit code policy TBD)

Avoid marketing output like `PARTIALLY ASSURED` in v0.1.

---

## 4. Compose, don’t reinvent

| Concern | Use |
|---|---|
| Component / model inventory | SPDX 3 / CycloneDX ML-BOM (external ref) |
| Security tool findings | SARIF (evidence type) |
| Build provenance | SLSA / in-toto / Sigstore (later) |
| Runtime telemetry | OTel GenAI pointers (URIs), not embedded traces |
| Regime classification | RiskLine (or peers) |
| EU article packs | Regula / EU profile (later interop) |

**Preferred serialization strategy (decision needed):**

- **A — Dual:** developer YAML sugar → optional export to CycloneDX Declarations.
- **B — CycloneDX-native:** AegisProof is a profile + CLI conventions on CDX (harder DX, better interop).
- **C — Independent schema only:** highest reinvention risk; only if CDX export is proven inadequate.

Recommendation: **A**, with explicit kill-switch to fold into **B** if dual maintenance hurts.

---

## 5. Integrity model

| Milestone | Integrity |
|---|---|
| v0.1 | SHA-256 digests + path confinement |
| v0.2 | Canonical JSON for signing; Ed25519 optional; git commit binding |
| Later | Sigstore / in-toto / RFC3161 timestamps |

Do not ship blockchain, transparency logs, or “immutable registries” in early milestones.

---

## 6. CLI surface (proposed)

```text
aegisproof init      # scaffold manifest + example evidence
aegisproof verify .  # default
aegisproof pack      # optional: hash-fill digests
aegisproof report    # alias of verify --format human
```

Outputs:

- human text
- `--format json`
- exit codes: `0` valid, `1` invalid, `2` stale-only (if enabled), `3` usage error

No LLM calls. No network required for verify (except optional future remote locators — off by default).

---

## 7. Local-first / cloud-optional

```text
Core: filesystem + CLI
Optional later:
  object storage evidence store
  audit API
  hosted verify
```

Cloud must remain an **optimisation**, not a dependency. Cost target: free tier / single maintainer affordable (static docs on GitHub Pages or similar; no always-on cluster).

---

## 8. Extension / profiles

```text
AegisProof Core (evidence package)
   ├── profile: eu-ai-act (maps evidence types → articles) — NOT in core
   ├── profile: nist-ai-rmf
   ├── profile: healthcare
   └── profile: financial-services
```

Profiles add **mapping documents** and conformance suites; they must not change core verify semantics.

---

## 9. Consumers

| Consumer | Role |
|---|---|
| CI | Gate on `INVALID` |
| RiskLine | Import package / emit package from assure |
| Auditors | Human-readable report + digests |
| Agents | Read-only retrieval via optional indexers |
| XERJ | Optional search index — **not** trust root |

---

## 10. Deliberate non-goals (v0.1)

- Risk scoring / legal classification
- Control implementation status theatre
- Mandatory cloud
- Mandatory LLM
- Mandatory XERJ
- Full CycloneDX reimplementation
- Enterprise dashboard
- “Assured / compliant” verdict language

---

## 11. Relationship to RiskLine architecture

RiskLine already implements classify → assure → signed evidence verify. AegisProof should **not** reimplement the regime engine.

Healthy split:

```text
RiskLine:  what controls/regime say about a system description
AegisProof: whether evidence artefacts for a system version are present + intact
```

Interop sketch: RiskLine can emit an AegisProof package from assure outputs; AegisProof never calls RiskLine.

---

## 12. Open architecture decisions

1. Dual YAML vs CycloneDX-native?
2. Include `claims[]` in v0.1 or defer?
3. Exit code for stale evidence: fail CI or warn?
4. Single-file package vs directory convention?
5. Extract RiskLine `EvidenceBundle` fields into AegisProof, or keep them RiskLine-private and only export?
