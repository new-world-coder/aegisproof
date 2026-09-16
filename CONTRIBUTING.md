# Contributing

Thanks for interest in AegisProof.

## Before coding

Read `SPEC.md` and `docs/ARCHITECTURE.md`. Core rule: **evidence ≠ compliance**.

## Development

```bash
go test ./...
go run ./cmd/aegisproof verify ./examples/minimal
```

## DCO

Commits should be signed off (`git commit -s`) under the Developer Certificate of Origin.

## Scope

- Spec/schema/CLI/conformance improvements welcome
- Do not add LLM calls, cloud requirements, or compliance verdict language to core verify
