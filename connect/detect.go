package connect

import (
	"os"
	"path/filepath"
	"sort"
)

func DetectPages(designsDir string, ignore []string) ([]string, error) {
	entries, err := os.ReadDir(designsDir)
	if err != nil {
		return nil, err
	}

	ignoreSet := make(map[string]bool)
	for _, name := range ignore {
		ignoreSet[name] = true
	}

	var pages []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if ignoreSet[name] {
			continue
		}
		indexPath := filepath.Join(designsDir, name, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			pages = append(pages, name)
		}
	}

	sort.Strings(pages)
	return pages, nil
}
