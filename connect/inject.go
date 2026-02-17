package connect

import (
	"errors"
	"os"
	"strings"
)

func InjectScript(htmlPath, script string, dryRun bool) error {
	data, err := os.ReadFile(htmlPath)
	if err != nil {
		return err
	}

	content := string(data)
	content = RemoveExistingScript(content)

	idx := findBodyCloseIndex(content)
	if idx == -1 {
		return errors.New("</body> tag not found")
	}

	injected := content[:idx] + script + "\n" + content[idx:]

	if !dryRun {
		if err := os.WriteFile(htmlPath, []byte(injected), 0644); err != nil {
			return err
		}
	}

	return nil
}

func findBodyCloseIndex(content string) int {
	lower := strings.ToLower(content)
	return strings.LastIndex(lower, "</body>")
}
