# Conformance fixtures (v0.1)

| Case | Path | Expect |
|---|---|---|
| valid-minimal | `valid/minimal` → `../../examples/minimal` | VALID |
| valid-rag | `valid/rag-claims` → `../../examples/rag-claims` | VALID |
| invalid-bad-digest | `invalid/bad-digest` | INVALID / DIGEST_MISMATCH |
| invalid-missing | `invalid/missing-file` | INVALID / EVIDENCE_MISSING |
| invalid-traversal | `invalid/path-traversal` | INVALID / PATH |
| invalid-schema | `invalid/bad-spec-version` | INVALID / SCHEMA |
| invalid-duplicate-id | `invalid/duplicate-id` | INVALID / SCHEMA |

Run:

```bash
go run ./cmd/aegisproof verify ./conformance/v0.1/valid/minimal
go run ./cmd/aegisproof verify ./conformance/v0.1/invalid/bad-digest; echo exit:$?
```
