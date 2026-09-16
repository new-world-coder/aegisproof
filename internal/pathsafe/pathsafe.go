package pathsafe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Resolve ensures locator stays within root after cleaning.
func Resolve(root, locator string) (string, error) {
	if locator == "" {
		return "", fmt.Errorf("empty locator")
	}
	if filepath.IsAbs(locator) {
		return "", fmt.Errorf("absolute locator not allowed: %s", locator)
	}
	// Normalise to slash then to OS for checks on ".."
	cleaned := filepath.Clean("/" + strings.ReplaceAll(locator, "\\", "/"))
	rel := strings.TrimPrefix(cleaned, "/")
	if rel == "" || rel == "." {
		return "", fmt.Errorf("invalid locator: %s", locator)
	}
	if strings.HasPrefix(rel, "..") || strings.Contains(rel, string(filepath.Separator)+"..") {
		return "", fmt.Errorf("path escape rejected: %s", locator)
	}
	// Also reject explicit .. segments in original
	for _, part := range strings.Split(strings.ReplaceAll(locator, "\\", "/"), "/") {
		if part == ".." {
			return "", fmt.Errorf("path escape rejected: %s", locator)
		}
	}

	full := filepath.Join(root, filepath.FromSlash(rel))
	fullClean := filepath.Clean(full)
	rootClean := filepath.Clean(root)

	relToRoot, err := filepath.Rel(rootClean, fullClean)
	if err != nil || strings.HasPrefix(relToRoot, "..") {
		return "", fmt.Errorf("path escape rejected: %s", locator)
	}

	// Symlink escape check
	resolved, err := filepath.EvalSymlinks(fullClean)
	if err != nil {
		if os.IsNotExist(err) {
			return fullClean, nil // missing handled by caller
		}
		// If parent doesn't exist, EvalSymlinks may fail — return cleaned path
		if _, statErr := os.Lstat(fullClean); statErr != nil {
			return fullClean, nil
		}
		return "", fmt.Errorf("resolve symlink: %w", err)
	}
	rootResolved, err := filepath.EvalSymlinks(rootClean)
	if err != nil {
		rootResolved = rootClean
	}
	relResolved, err := filepath.Rel(rootResolved, resolved)
	if err != nil || strings.HasPrefix(relResolved, "..") {
		return "", fmt.Errorf("symlink escape rejected: %s", locator)
	}
	return resolved, nil
}
