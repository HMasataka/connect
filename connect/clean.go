package connect

import (
	"os"
	"regexp"
)

var scriptPattern = regexp.MustCompile(`(?s)<script data-connect>.*?</script>\n?`)

func RemoveExistingScript(content string) string {
	return scriptPattern.ReplaceAllString(content, "")
}

func CleanScript(htmlPath string, dryRun bool) (bool, error) {
	data, err := os.ReadFile(htmlPath)
	if err != nil {
		return false, err
	}

	content := string(data)
	cleaned := RemoveExistingScript(content)

	if content == cleaned {
		return false, nil
	}

	if !dryRun {
		if err := os.WriteFile(htmlPath, []byte(cleaned), 0644); err != nil {
			return false, err
		}
	}

	return true, nil
}
