# Quick start

```bash
git clone https://github.com/new-world-coder/aegisproof.git
cd aegisproof
go test ./...
go run ./cmd/aegisproof verify ./examples/minimal
go run ./cmd/aegisproof init ./my-system
go run ./cmd/aegisproof pack ./my-system
go run ./cmd/aegisproof verify ./my-system
```

Requires Go 1.22+.
