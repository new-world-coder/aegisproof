package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const SpecVersion = "0.1"

type Package struct {
	SpecVersion       string             `json:"specVersion" yaml:"specVersion"`
	System            System             `json:"system" yaml:"system"`
	Producer          Producer           `json:"producer" yaml:"producer"`
	CreatedAt         string             `json:"createdAt" yaml:"createdAt"`
	Description       string             `json:"description,omitempty" yaml:"description,omitempty"`
	Profiles          []string           `json:"profiles,omitempty" yaml:"profiles,omitempty"`
	ExternalDocuments []ExternalDocument `json:"externalDocuments,omitempty" yaml:"externalDocuments,omitempty"`
	Evidence          []EvidenceItem     `json:"evidence" yaml:"evidence"`
}

type System struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	ID      string `json:"id,omitempty" yaml:"id,omitempty"`
}

type Producer struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

type ExternalDocument struct {
	Type    string `json:"type" yaml:"type"`
	Locator string `json:"locator" yaml:"locator"`
}

type EvidenceItem struct {
	ID         string `json:"id" yaml:"id"`
	Type       string `json:"type" yaml:"type"`
	Locator    string `json:"locator" yaml:"locator"`
	Digest     string `json:"digest" yaml:"digest"`
	MediaType  string `json:"mediaType,omitempty" yaml:"mediaType,omitempty"`
	ProducedAt string `json:"producedAt,omitempty" yaml:"producedAt,omitempty"`
	ExpiresAt  string `json:"expiresAt,omitempty" yaml:"expiresAt,omitempty"`
}

type Loaded struct {
	Root     string
	Path     string
	Format   string // "yaml" or "json"
	Package  Package
	RawBytes []byte
}

func FindManifest(root string) (string, error) {
	jsonPath := filepath.Join(root, "aegisproof.json")
	yamlPath := filepath.Join(root, "aegisproof.yaml")
	ymlPath := filepath.Join(root, "aegisproof.yml")

	if st, err := os.Stat(jsonPath); err == nil && !st.IsDir() {
		return jsonPath, nil
	}
	if st, err := os.Stat(yamlPath); err == nil && !st.IsDir() {
		return yamlPath, nil
	}
	if st, err := os.Stat(ymlPath); err == nil && !st.IsDir() {
		return ymlPath, nil
	}
	return "", fmt.Errorf("no aegisproof.json or aegisproof.yaml found in %s", root)
}

func Load(root string) (*Loaded, error) {
	root = filepath.Clean(root)
	path, err := FindManifest(root)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	const maxManifest = 1 << 20
	if len(raw) > maxManifest {
		return nil, fmt.Errorf("manifest exceeds 1 MiB limit")
	}

	var pkg Package
	format := "json"
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		format = "yaml"
		if err := yaml.Unmarshal(raw, &pkg); err != nil {
			return nil, fmt.Errorf("parse YAML manifest: %w", err)
		}
	default:
		if err := json.Unmarshal(raw, &pkg); err != nil {
			return nil, fmt.Errorf("parse JSON manifest: %w", err)
		}
	}

	return &Loaded{Root: root, Path: path, Format: format, Package: pkg, RawBytes: raw}, nil
}

func (p Package) StructuralValidate() []string {
	var errs []string
	if p.SpecVersion != SpecVersion {
		errs = append(errs, fmt.Sprintf("unsupported specVersion %q (want %q)", p.SpecVersion, SpecVersion))
	}
	if strings.TrimSpace(p.System.Name) == "" {
		errs = append(errs, "system.name is required")
	}
	if strings.TrimSpace(p.System.Version) == "" {
		errs = append(errs, "system.version is required")
	}
	if strings.TrimSpace(p.Producer.Name) == "" {
		errs = append(errs, "producer.name is required")
	}
	if _, err := time.Parse(time.RFC3339, p.CreatedAt); err != nil {
		errs = append(errs, "createdAt must be RFC3339")
	}
	if len(p.Evidence) == 0 {
		errs = append(errs, "evidence must contain at least one item")
	}
	seen := map[string]struct{}{}
	for i, e := range p.Evidence {
		if e.ID == "" {
			errs = append(errs, fmt.Sprintf("evidence[%d].id is required", i))
		} else if _, ok := seen[e.ID]; ok {
			errs = append(errs, fmt.Sprintf("duplicate evidence id %q", e.ID))
		} else {
			seen[e.ID] = struct{}{}
		}
		if !ValidEvidenceType(e.Type) {
			errs = append(errs, fmt.Sprintf("evidence[%d].type %q is invalid", i, e.Type))
		}
		if e.Locator == "" {
			errs = append(errs, fmt.Sprintf("evidence[%d].locator is required", i))
		}
		if !ValidDigest(e.Digest) {
			errs = append(errs, fmt.Sprintf("evidence[%d].digest must match sha256:[64 hex]", i))
		}
		if e.ProducedAt != "" {
			if _, err := time.Parse(time.RFC3339, e.ProducedAt); err != nil {
				errs = append(errs, fmt.Sprintf("evidence[%d].producedAt must be RFC3339", i))
			}
		}
		if e.ExpiresAt != "" {
			if _, err := time.Parse(time.RFC3339, e.ExpiresAt); err != nil {
				errs = append(errs, fmt.Sprintf("evidence[%d].expiresAt must be RFC3339", i))
			}
		}
	}
	return errs
}

func ValidEvidenceType(t string) bool {
	switch t {
	case "evaluation", "policy", "test-result", "sbom", "sarif", "model-card", "runtime-receipt", "other":
		return true
	default:
		return false
	}
}

func ValidDigest(d string) bool {
	if !strings.HasPrefix(d, "sha256:") {
		return false
	}
	hexPart := d[len("sha256:"):]
	if len(hexPart) != 64 {
		return false
	}
	for _, c := range hexPart {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
