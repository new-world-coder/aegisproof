package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const MaxEvidenceBytes = 64 << 20 // 64 MiB

func FileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return "", err
	}
	if st.Size() > MaxEvidenceBytes {
		return "", fmt.Errorf("file exceeds 64 MiB limit: %s", path)
	}

	h := sha256.New()
	limited := io.LimitReader(f, MaxEvidenceBytes+1)
	n, err := io.Copy(h, limited)
	if err != nil {
		return "", err
	}
	if n > MaxEvidenceBytes {
		return "", fmt.Errorf("file exceeds 64 MiB limit: %s", path)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}
