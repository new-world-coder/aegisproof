# AegisProof ↔ XERJ Integration Hypothesis (Phase 0)

**Status:** Exploratory technical hypothesis only.  
**Not a partnership, endorsement, commitment, or dependency.**

XERJ: https://xerj.org/ · https://github.com/xerj-org/xerj

---

## 1. What XERJ is (technical)

From public materials:

- Single Rust binary, Apache-2.0, Elasticsearch-compatible HTTP API.
- **`xerj autoindex`** — content-sniffed ingest across many formats (including YAML/JSON).
- Tree-sitter based code indexing for reference-coding agents.
- `/_memory` namespaced agent memory; `/_graph` knowledge layer with evidence on links.
- Value prop: reduce agent token burn by retrieving real implementations vs inventing them.

Measured case study focus: **reference-coding**, not assurance evidence integrity.

---

## 2. Fit assessment for AegisProof

| Question | Assessment |
|---|---|
| Index AegisProof YAML/JSON efficiently? | **Yes, likely** — autoindex already handles YAML/JSON; small manifests are trivial. |
| Preserve provenance when indexing? | **Partial** — index is a search projection. Digests/URIs must remain in source artefacts; XERJ must not become canonical store. |
| Retrieve without mutating evidence? | **Yes, if used read-only** — search/recall must not rewrite package files. |
| Point back to immutable evidence? | **Yes, if designed** — store `digest`, `locator`, `package_id` as fields; retrieval returns pointers, not “truth.” |
| Evidence knowledge layer? | **Possible demo** — organisational search over many packages. Not required for standard success. |
| Discover reference implementations for controls? | **Strong XERJ native fit** — orthogonal to AegisProof; useful for *implementing* controls, not proving them. |

**Strategic centrality:** **Nice-to-have / networking / later experiment** — not on the critical path for a standard. Treating XERJ as architectural pillar would confuse product identity (search infra vs evidence format).

---

## 3. Recommended stance for founder

1. Keep conversations with Ivan **exploratory and technical**.
2. Do not put XERJ in AegisProof core README architecture diagrams as a required box.
3. After v0.1 package exists, run cheap indexing experiments.
4. Stronger near-term story: XERJ helps agents **find** control implementations; AegisProof helps CI **verify** evidence packs — complementary layers, optional glue.

```text
                 AI AGENT
                    │
          ┌─────────┴──────────┐
          ▼                    ▼
       XERJ                 AegisProof
   retrieval (optional)    evidence (source of integrity)
          │                    │
          └─────────┬──────────┘
                    ▼
            optional RiskLine
```

This is a **research hypothesis**, not a final architecture.

---

## 4. Four cheap experiments

### Experiment 1 — Index packages

- Author 5 sample AegisProof packages (fixture).
- `xerj autoindex` the fixtures directory.
- Query: “PII protection evaluation evidence for claims-assistant”.
- **Pass:** correct package/file in top-3; digests visible in hit source.

### Experiment 2 — Agent retrieve, don’t mutate

- Agent task: “List evidence digests for system X.”
- Agent may only use XERJ search + read files.
- **Pass:** answers match `aegisproof verify` JSON; working tree unchanged.

### Experiment 3 — Reference coding + evidence

- Agent implements a logging control using XERJ reference-coding against a known OSS pattern.
- Separately attaches resulting test output as AegisProof evidence item.
- **Pass:** both retrieval and verify succeed; licences respected.

### Experiment 4 — Workflow timebox

- Two arms: (A) find evidence via `find`/`rg`, (B) via XERJ.
- Measure minutes to assemble questionnaire answers for a fixture app.
- **Pass:** ≥20% time reduction **or** document null result honestly (XERJ may not help tiny repos).

---

## 5. What must NOT be coupled

- XERJ as required runtime for `aegisproof verify`
- Storing only indexed copies without file digests
- Letting agent memory overwrite evidence
- Marketing “XERJ + RiskLine + AegisProof partnership”
- Embedding Elasticsearch-js contribution as product dependency (personal credibility ≠ architecture)

---

## 6. Priority vs core work

| Priority | Work |
|---|---|
| P0 | Spec + examples + verify CLI |
| P1 | CI action + RiskLine interop experiment |
| P2 | XERJ experiments (above) |
| P3 | Shared blog/case study *if* experiments produce honest numbers |

---

## 7. Collaboration proposal sketch (for later human use)

Tone: technical experiment invitation, not commercial MoU.

> We are designing a regulation-neutral AI evidence package format. Separately, XERJ’s autoindex may help agents discover evidence and reference implementations. We propose a time-boxed indexing experiment with sample packages — no dependency either way. Happy to share fixtures and measure retrieval quality together.
