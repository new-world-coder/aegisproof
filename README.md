# AegisProof

**Open, vendor-neutral evidence packages for AI systems.**

AegisProof is a small specification and CLI so teams can declare *what evidence exists* for an AI system revision (`aegisproof.yaml`), hash those artefacts, and verify integrity in CI — without an LLM, without a SaaS, and without claiming legal compliance.

> **Integrity ≠ compliance.** A green `verify` means listed files match their digests. It does **not** mean the system is safe, certified, EU AI Act compliant, or “assured.”

If you ship LLM/RAG/agent features into enterprises, you already assemble evals, policies, test outputs, model cards, and SBOMs for security questionnaires. AegisProof turns that pile into a **portable, diffable package** other tools (and humans) can check.

---

## Why fork / contribute

- **Standards-shaped OSS** — schema + conformance + boring CLI (Apache-2.0)
- **Compose, don’t reinvent** — point at SPDX / CycloneDX / SARIF; don’t replace them
- **Local-first** — works offline; cloud is optional and not required
- **Clear non-goals** — not a GRC dashboard, not RiskLine, not “AI SBOM”

Good first moves: run the examples, add a conformance fixture, emit a package from your AI app’s CI, open an issue.

**Docs / wiki:** [`docs/wiki/Home.md`](docs/wiki/Home.md) · **Issues:** [open issues](https://github.com/new-world-coder/aegisproof/issues) · **Spec:** [`SPEC.md`](SPEC.md)

---

## Quick start

Requires [Go](https://go.dev/dl/) 1.22+.

```bash
git clone https://github.com/new-world-coder/aegisproof.git
cd aegisproof

go test ./...
go run ./cmd/aegisproof verify ./examples/minimal
go run ./cmd/aegisproof verify ./examples/minimal --format json
```

Scaffold your own package:

```bash
go run ./cmd/aegisproof init ./my-system
# replace evidence files with real artefacts…
go run ./cmd/aegisproof pack ./my-system
go run ./cmd/aegisproof verify ./my-system
```

### Minimal manifest shape

```yaml
specVersion: "0.1"
system:
  name: my-ai-system
  version: "1.0.0"
producer:
  name: my-team
createdAt: "2026-09-16T12:00:00Z"
evidence:
  - id: eval-results
    type: evaluation
    locator: evidence/eval.json
    digest: "sha256:…"   # filled by `aegisproof pack`
```

---

## What this is / is not

| Is | Is not |
|---|---|
| Evidence **package** + digest verify | Legal compliance engine |
| CI-friendly interchange (SARIF-like role) | “AI SBOM” (use SPDX/CycloneDX for inventory) |
| Regulation-**neutral** core | EU AI Act / NIST encoded in core schema |
| Human-readable YAML/JSON | LLM-required tooling |
| Optional consumer for tools like RiskLine | A RiskLine feature / proprietary format |

---

## Repository map

| Path | Purpose |
|---|---|
| [`SPEC.md`](SPEC.md) | v0.1 specification |
| [`schema/`](schema/) | JSON Schema |
| [`cmd/aegisproof`](cmd/aegisproof) | CLI (`verify`, `pack`, `init`) |
| [`examples/`](examples/) | minimal, RAG, agent packages |
| [`conformance/`](conformance/) | valid / invalid fixtures |
| [`docs/`](docs/) | architecture, research, roadmap, security |

### Exit codes

| Code | Result |
|---|---|
| 0 | `VALID` |
| 1 | `INVALID` |
| 2 | `STALE` (digests OK; freshness warnings) |
| 3 | Usage / IO error |

---

## Related projects (not dependencies)

- **[RiskLine](https://github.com/new-world-coder/riskline)** — deterministic AI risk classification / assure (may consume AegisProof later; does not own this standard)
- **SPDX / CycloneDX** — inventory & attestations to *reference*, not replace
- **XERJ** — optional retrieval/indexing experiments only

---

## Status

Early **v0.1 draft**. APIs and schema may change before 1.0. Feedback via [Issues](https://github.com/new-world-coder/aegisproof/issues) is welcome — especially emitter/consumer designs and conformance gaps.

Planning background: [`docs/IMPLEMENTATION_AND_TESTING_STRATEGY.md`](docs/IMPLEMENTATION_AND_TESTING_STRATEGY.md).

## Licence

Apache-2.0 — see [`LICENSE`](LICENSE).
