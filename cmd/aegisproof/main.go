package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aegisproof/aegisproof/internal/digest"
	"github.com/aegisproof/aegisproof/internal/manifest"
	"github.com/aegisproof/aegisproof/internal/verify"
	"gopkg.in/yaml.v3"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(3)
	}
	switch os.Args[1] {
	case "verify":
		os.Exit(cmdVerify(os.Args[2:]))
	case "pack":
		os.Exit(cmdPack(os.Args[2:]))
	case "init":
		os.Exit(cmdInit(os.Args[2:]))
	case "version", "--version", "-v":
		fmt.Printf("aegisproof %s\n", version)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		usage()
		os.Exit(3)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `AegisProof — AI assurance evidence package verifier

Usage:
  aegisproof verify [dir] [--format human|json] [--strict-stale]
  aegisproof pack [dir]          Fill/update digests in the manifest
  aegisproof init [dir]          Scaffold a minimal package
  aegisproof version

Exit codes (verify): 0 VALID, 1 INVALID, 2 STALE, 3 usage/IO

Disclaimer: verify checks package integrity only — not legal compliance.
`)
}

func cmdVerify(args []string) int {
	dir := "."
	format := "human"
	strictStale := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--format requires a value")
				return 3
			}
			i++
			format = args[i]
		case "--strict-stale":
			strictStale = true
		case "-h", "--help":
			usage()
			return 0
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(os.Stderr, "unknown flag %s\n", args[i])
				return 3
			}
			dir = args[i]
		}
	}

	rep, err := verify.Package(dir, verify.Options{StrictStale: strictStale})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rep)
	case "human":
		printHuman(rep)
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q\n", format)
		return 3
	}
	return verify.ExitCode(rep.Result)
}

func printHuman(rep *verify.Report) {
	fmt.Printf("AegisProof Integrity Report\n\n")
	fmt.Printf("System:   %s %s\n", rep.SystemName, rep.SystemVersion)
	fmt.Printf("Manifest: %s\n", rep.Manifest)
	fmt.Printf("Evidence: %d/%d digests matched\n", rep.EvidenceOK, rep.EvidenceTotal)
	fmt.Printf("Result:   %s\n", rep.Result)
	if len(rep.Errors) > 0 {
		fmt.Printf("\nErrors\n")
		for _, e := range rep.Errors {
			if e.EvidenceID != "" {
				fmt.Printf("  ✗ [%s] %s (%s)\n", e.Code, e.Message, e.EvidenceID)
			} else {
				fmt.Printf("  ✗ [%s] %s\n", e.Code, e.Message)
			}
		}
	}
	if len(rep.Warnings) > 0 {
		fmt.Printf("\nWarnings\n")
		for _, w := range rep.Warnings {
			if w.EvidenceID != "" {
				fmt.Printf("  ⚠ [%s] %s (%s)\n", w.Code, w.Message, w.EvidenceID)
			} else {
				fmt.Printf("  ⚠ [%s] %s\n", w.Code, w.Message)
			}
		}
	}
	fmt.Printf("\n%s\n", rep.Disclaimer)
}

func cmdPack(args []string) int {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	loaded, err := manifest.Load(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}
	pkg := loaded.Package
	for i := range pkg.Evidence {
		item := &pkg.Evidence[i]
		full := filepath.Join(loaded.Root, filepath.FromSlash(item.Locator))
		sum, err := digest.FileSHA256(full)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error hashing %s: %v\n", item.Locator, err)
			return 3
		}
		item.Digest = sum
	}
	out, err := yaml.Marshal(&pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}
	// Prefer writing yaml path
	outPath := loaded.Path
	if loaded.Format == "json" {
		outPath = filepath.Join(loaded.Root, "aegisproof.json")
		jb, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 3
		}
		jb = append(jb, '\n')
		if err := os.WriteFile(outPath, jb, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 3
		}
	} else {
		if err := os.WriteFile(outPath, out, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 3
		}
	}
	fmt.Printf("updated digests in %s\n", outPath)
	return 0
}

func cmdInit(args []string) int {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	if err := os.MkdirAll(filepath.Join(dir, "evidence"), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}
	sample := filepath.Join(dir, "evidence", "placeholder.txt")
	content := []byte("Replace this file with real evidence artefacts.\n")
	if _, err := os.Stat(sample); os.IsNotExist(err) {
		if err := os.WriteFile(sample, content, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 3
		}
	}
	sum, err := digest.FileSHA256(sample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}
	pkg := manifest.Package{
		SpecVersion: manifest.SpecVersion,
		System:      manifest.System{Name: "my-ai-system", Version: "0.1.0"},
		Producer:    manifest.Producer{Name: "aegisproof-init", Version: version},
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		Description: "Scaffolded by aegisproof init. Integrity only — not a compliance verdict.",
		Evidence: []manifest.EvidenceItem{{
			ID:        "placeholder",
			Type:      "other",
			Locator:   "evidence/placeholder.txt",
			Digest:    sum,
			MediaType: "text/plain",
		}},
	}
	outPath := filepath.Join(dir, "aegisproof.yaml")
	if _, err := os.Stat(outPath); err == nil {
		fmt.Fprintf(os.Stderr, "error: %s already exists\n", outPath)
		return 3
	}
	out, err := yaml.Marshal(&pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}
	fmt.Printf("created %s\n", outPath)
	fmt.Printf("next: replace evidence files, run `aegisproof pack %s`, then `aegisproof verify %s`\n", dir, dir)
	return 0
}
