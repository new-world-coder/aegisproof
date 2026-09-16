package verify

import (
	"fmt"
	"os"
	"time"

	"github.com/aegisproof/aegisproof/internal/digest"
	"github.com/aegisproof/aegisproof/internal/manifest"
	"github.com/aegisproof/aegisproof/internal/pathsafe"
)

type Result string

const (
	ResultValid   Result = "VALID"
	ResultStale   Result = "STALE"
	ResultInvalid Result = "INVALID"
)

type Issue struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	EvidenceID string `json:"evidenceId,omitempty"`
}

type Report struct {
	Result      Result  `json:"result"`
	SpecVersion string  `json:"specVersion"`
	SystemName  string  `json:"systemName"`
	SystemVersion string `json:"systemVersion"`
	Manifest    string  `json:"manifest"`
	CheckedAt   string  `json:"checkedAt"`
	EvidenceOK  int     `json:"evidenceOk"`
	EvidenceTotal int   `json:"evidenceTotal"`
	Errors      []Issue `json:"errors"`
	Warnings    []Issue `json:"warnings"`
	Disclaimer  string  `json:"disclaimer"`
}

const Disclaimer = "Integrity check only. VALID/STALE does not mean compliant, certified, assured, or safe."

type Options struct {
	StrictStale bool
	Now         time.Time
}

func Package(root string, opts Options) (*Report, error) {
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	loaded, err := manifest.Load(root)
	if err != nil {
		return nil, err
	}
	pkg := loaded.Package
	rep := &Report{
		Result:        ResultValid,
		SpecVersion:   pkg.SpecVersion,
		SystemName:    pkg.System.Name,
		SystemVersion: pkg.System.Version,
		Manifest:      loaded.Path,
		CheckedAt:     opts.Now.Format(time.RFC3339),
		EvidenceTotal: len(pkg.Evidence),
		Disclaimer:    Disclaimer,
	}

	for _, msg := range pkg.StructuralValidate() {
		rep.Errors = append(rep.Errors, Issue{Code: "SCHEMA", Message: msg})
	}
	if len(rep.Errors) > 0 {
		rep.Result = ResultInvalid
		return rep, nil
	}

	if len(pkg.Evidence) > 10000 {
		rep.Errors = append(rep.Errors, Issue{Code: "LIMIT", Message: "evidence count exceeds 10000"})
		rep.Result = ResultInvalid
		return rep, nil
	}

	for _, item := range pkg.Evidence {
		full, err := pathsafe.Resolve(loaded.Root, item.Locator)
		if err != nil {
			rep.Errors = append(rep.Errors, Issue{Code: "PATH", Message: err.Error(), EvidenceID: item.ID})
			continue
		}
		st, err := os.Stat(full)
		if err != nil {
			if os.IsNotExist(err) {
				rep.Errors = append(rep.Errors, Issue{Code: "EVIDENCE_MISSING", Message: fmt.Sprintf("missing file %s", item.Locator), EvidenceID: item.ID})
			} else {
				rep.Errors = append(rep.Errors, Issue{Code: "IO", Message: err.Error(), EvidenceID: item.ID})
			}
			continue
		}
		if st.IsDir() {
			rep.Errors = append(rep.Errors, Issue{Code: "PATH", Message: "locator points to directory", EvidenceID: item.ID})
			continue
		}
		got, err := digest.FileSHA256(full)
		if err != nil {
			rep.Errors = append(rep.Errors, Issue{Code: "IO", Message: err.Error(), EvidenceID: item.ID})
			continue
		}
		if got != item.Digest {
			rep.Errors = append(rep.Errors, Issue{
				Code:       "DIGEST_MISMATCH",
				Message:    fmt.Sprintf("expected %s, got %s", item.Digest, got),
				EvidenceID: item.ID,
			})
			continue
		}
		rep.EvidenceOK++

		if item.ExpiresAt != "" {
			exp, err := time.Parse(time.RFC3339, item.ExpiresAt)
			if err == nil && opts.Now.After(exp) {
				issue := Issue{Code: "STALE", Message: fmt.Sprintf("expiresAt %s is in the past", item.ExpiresAt), EvidenceID: item.ID}
				if opts.StrictStale {
					rep.Errors = append(rep.Errors, issue)
				} else {
					rep.Warnings = append(rep.Warnings, issue)
				}
			}
		}
	}

	if len(rep.Errors) > 0 {
		rep.Result = ResultInvalid
		return rep, nil
	}
	if len(rep.Warnings) > 0 {
		rep.Result = ResultStale
		return rep, nil
	}
	rep.Result = ResultValid
	return rep, nil
}

func ExitCode(r Result) int {
	switch r {
	case ResultValid:
		return 0
	case ResultInvalid:
		return 1
	case ResultStale:
		return 2
	default:
		return 3
	}
}
