package connect

import (
	"strings"
	"testing"
)

func TestGenerateScript_EmptyConfig(t *testing.T) {
	// Given
	pages := []string{"home", "about"}
	cfg := Config{
		Links:  []LinkRule{},
		Modals: []string{},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `var pages = ["home","about"]`) {
		t.Errorf("expected pages array in script, got:\n%s", result)
	}
	if !strings.Contains(result, `var modals = []`) {
		t.Errorf("expected empty modals array in script, got:\n%s", result)
	}
	if !strings.Contains(result, `var links = []`) {
		t.Errorf("expected empty links array in script, got:\n%s", result)
	}
	if !strings.Contains(result, `var close = []`) {
		t.Errorf("expected empty close array in script, got:\n%s", result)
	}
	if !strings.Contains(result, `<script data-connect>`) {
		t.Errorf("expected script tag in output, got:\n%s", result)
	}
}

func TestGenerateScript_WithLinks(t *testing.T) {
	// Given
	pages := []string{"home", "about", "settings"}
	cfg := Config{
		Links: []LinkRule{
			{
				Selector: ".nav-item",
				Match:    "text",
				Mapping:  map[string]string{"Home": "home"},
			},
		},
		Modals: []string{},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `"selector":".nav-item"`) {
		t.Errorf("expected selector in links, got:\n%s", result)
	}
	if !strings.Contains(result, `"match":"text"`) {
		t.Errorf("expected match in links, got:\n%s", result)
	}
	if !strings.Contains(result, `"mapping":{"Home":"home"}`) {
		t.Errorf("expected mapping in links, got:\n%s", result)
	}
}

func TestGenerateScript_WithTarget(t *testing.T) {
	// Given
	pages := []string{"home", "branches-dialog"}
	cfg := Config{
		Links: []LinkRule{
			{
				Selector: ".branch-selector",
				Match:    "text",
				Target:   "branches-dialog",
			},
		},
		Modals: []string{},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `"target":"branches-dialog"`) {
		t.Errorf("expected target in links, got:\n%s", result)
	}
}

func TestGenerateScript_WithModals(t *testing.T) {
	// Given
	pages := []string{"home", "tags", "settings"}
	cfg := Config{
		Links:  []LinkRule{},
		Modals: []string{"tags"},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `var modals = ["tags"]`) {
		t.Errorf("expected modals array in script, got:\n%s", result)
	}
}

func TestGenerateScript_WithCloseRules(t *testing.T) {
	// Given
	pages := []string{"home"}
	cfg := Config{
		Links:  []LinkRule{},
		Modals: []string{},
		Close: []CloseRule{
			{Selector: ".modal-close", Match: "", Value: ""},
			{Selector: ".btn", Match: "text", Value: "Close"},
		},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `"selector":".modal-close"`) {
		t.Errorf("expected close selector in script, got:\n%s", result)
	}
	if !strings.Contains(result, `"selector":".btn"`) {
		t.Errorf("expected close selector .btn in script, got:\n%s", result)
	}
	if !strings.Contains(result, `"match":"text"`) {
		t.Errorf("expected close match 'text' in script, got:\n%s", result)
	}
	if !strings.Contains(result, `"value":"Close"`) {
		t.Errorf("expected close value 'Close' in script, got:\n%s", result)
	}
}

func TestGenerateScript_ScriptStructure(t *testing.T) {
	// Given
	pages := []string{"home"}
	cfg := Config{
		Links:  []LinkRule{},
		Modals: []string{},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then: verify essential script functions and logic exist
	requiredSnippets := []string{
		"function toPageId(text)",
		"function getText(el, match)",
		"function resolve(text, mapping, target)",
		"function navigate(pageId)",
		"function closeModal()",
		"sessionStorage.setItem('connect-last-page', currentPage)",
		"sessionStorage.getItem('connect-last-page')",
		"links.forEach",
		"close.forEach",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(result, snippet) {
			t.Errorf("expected '%s' in script, got:\n%s", snippet, result)
		}
	}
}

func TestGenerateScript_NullTargetForEmptyString(t *testing.T) {
	// Given
	pages := []string{"home"}
	cfg := Config{
		Links: []LinkRule{
			{
				Selector: ".nav-item",
				Match:    "text",
				Target:   "",
			},
		},
		Modals: []string{},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `"target":null`) {
		t.Errorf("expected null target when empty string, got:\n%s", result)
	}
}

func TestGenerateScript_NullMatchAndValueForEmptyCloseRule(t *testing.T) {
	// Given
	pages := []string{"home"}
	cfg := Config{
		Links:  []LinkRule{},
		Modals: []string{},
		Close: []CloseRule{
			{Selector: ".modal-close", Match: "", Value: ""},
		},
		Ignore: []string{},
	}

	// When
	result := GenerateScript(pages, cfg)

	// Then
	if !strings.Contains(result, `"match":null`) {
		t.Errorf("expected null match when empty string, got:\n%s", result)
	}
	if !strings.Contains(result, `"value":null`) {
		t.Errorf("expected null value when empty string, got:\n%s", result)
	}
}
