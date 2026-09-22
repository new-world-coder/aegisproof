package verify_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aegisproof/aegisproof/internal/digest"
	"github.com/aegisproof/aegisproof/internal/verify"
)

func writePkg(t *testing.T, root string, body string, files map[string]string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "aegisproof.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	for rel, content := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestVerifyValid(t *testing.T) {
	root := t.TempDir()
	content := "hello evidence\n"
	sum, err := writeHashed(t, root, "evidence/a.txt", content)
	if err != nil {
		t.Fatal(err)
	}
	writePkg(t, root, `
specVersion: "0.1"
system: { name: t, version: "1" }
producer: { name: test }
createdAt: "2026-09-16T12:00:00Z"
evidence:
  - id: a
    type: other
    locator: evidence/a.txt
    digest: "`+sum+`"
`, nil)
	rep, err := verify.Package(root, verify.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Result != verify.ResultValid {
		t.Fatalf("got %s errors=%v", rep.Result, rep.Errors)
	}
}

func TestDigestMismatch(t *testing.T) {
	root := t.TempDir()
	writePkg(t, root, `
specVersion: "0.1"
system: { name: t, version: "1" }
producer: { name: test }
createdAt: "2026-09-16T12:00:00Z"
evidence:
  - id: a
    type: other
    locator: evidence/a.txt
    digest: "sha256:0000000000000000000000000000000000000000000000000000000000000000"
`, map[string]string{"evidence/a.txt": "x"})
	rep, err := verify.Package(root, verify.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Result != verify.ResultInvalid {
		t.Fatalf("want INVALID got %s", rep.Result)
	}
}

func TestDuplicateEvidenceID(t *testing.T) {
	root := t.TempDir()
	sum, err := writeHashed(t, root, "evidence/a.txt", "fixture\n")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "evidence", "b.txt"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writePkg(t, root, `
specVersion: "0.1"
system: { name: t, version: "1" }
producer: { name: test }
createdAt: "2026-09-16T12:00:00Z"
evidence:
  - id: a
    type: other
    locator: evidence/a.txt
    digest: "`+sum+`"
  - id: a
    type: other
    locator: evidence/b.txt
    digest: "`+sum+`"
`, nil)
	rep, err := verify.Package(root, verify.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Result != verify.ResultInvalid {
		t.Fatalf("want INVALID got %s errors=%v", rep.Result, rep.Errors)
	}
	if verify.ExitCode(rep.Result) != 1 {
		t.Fatalf("want exit 1 got %d", verify.ExitCode(rep.Result))
	}
	found := false
	for _, e := range rep.Errors {
		if e.Code == "SCHEMA" && strings.Contains(e.Message, `duplicate evidence id "a"`) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("want SCHEMA error for duplicate id, got %#v", rep.Errors)
	}
}

func TestConformanceInvalidDuplicateID(t *testing.T) {
	root := filepath.Join("..", "..", "conformance", "v0.1", "invalid", "duplicate-id")
	if _, err := os.Stat(filepath.Join(root, "aegisproof.yaml")); err != nil {
		t.Skip("conformance fixture not present:", err)
	}
	rep, err := verify.Package(root, verify.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Result != verify.ResultInvalid {
		t.Fatalf("want INVALID got %s errors=%v", rep.Result, rep.Errors)
	}
	if verify.ExitCode(rep.Result) != 1 {
		t.Fatalf("want exit 1 got %d", verify.ExitCode(rep.Result))
	}
}

func TestPathTraversal(t *testing.T) {
	root := t.TempDir()
	writePkg(t, root, `
specVersion: "0.1"
system: { name: t, version: "1" }
producer: { name: test }
createdAt: "2026-09-16T12:00:00Z"
evidence:
  - id: a
    type: other
    locator: ../outside.txt
    digest: "sha256:0000000000000000000000000000000000000000000000000000000000000000"
`, nil)
	rep, err := verify.Package(root, verify.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Result != verify.ResultInvalid {
		t.Fatalf("want INVALID got %s", rep.Result)
	}
}

func TestStaleWarning(t *testing.T) {
	root := t.TempDir()
	sum, err := writeHashed(t, root, "evidence/a.txt", "stale\n")
	if err != nil {
		t.Fatal(err)
	}
	writePkg(t, root, `
specVersion: "0.1"
system: { name: t, version: "1" }
producer: { name: test }
createdAt: "2026-09-16T12:00:00Z"
evidence:
  - id: a
    type: other
    locator: evidence/a.txt
    digest: "`+sum+`"
    expiresAt: "2020-01-01T00:00:00Z"
`, nil)
	rep, err := verify.Package(root, verify.Options{Now: time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Result != verify.ResultStale {
		t.Fatalf("want STALE got %s %#v", rep.Result, rep)
	}
}

func writeHashed(t *testing.T, root, rel, content string) (string, error) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		return "", err
	}
	return digest.FileSHA256(p)
}
