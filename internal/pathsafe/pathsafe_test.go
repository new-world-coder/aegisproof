package pathsafe_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aegisproof/aegisproof/internal/pathsafe"
)

func TestResolveRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	for _, loc := range []string{"../secret", "..\\secret", "/etc/passwd", "a/../../etc/passwd"} {
		if _, err := pathsafe.Resolve(root, loc); err == nil {
			t.Fatalf("expected reject for %q", loc)
		}
	}
}

func TestResolveAcceptsNested(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "evidence", "a.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := pathsafe.Resolve(root, "evidence/a.json")
	if err != nil {
		t.Fatal(err)
	}
	if got != path && filepath.Clean(got) != filepath.Clean(path) {
		// EvalSymlinks may change path form; ensure same file
		gi, _ := os.Stat(got)
		pi, _ := os.Stat(path)
		if !os.SameFile(gi, pi) {
			t.Fatalf("got %s want %s", got, path)
		}
	}
}
