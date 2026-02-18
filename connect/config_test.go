package connect

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCloseRule_UnmarshalJSON_String(t *testing.T) {
	// Given
	input := `".modal-close"`

	// When
	var rule CloseRule
	err := json.Unmarshal([]byte(input), &rule)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Selector != ".modal-close" {
		t.Errorf("expected selector '.modal-close', got '%s'", rule.Selector)
	}
	if rule.Match != "" {
		t.Errorf("expected match '', got '%s'", rule.Match)
	}
	if rule.Value != "" {
		t.Errorf("expected value '', got '%s'", rule.Value)
	}
}

func TestCloseRule_UnmarshalJSON_Object(t *testing.T) {
	// Given
	input := `{"selector": ".btn-secondary", "match": "text", "value": "Close"}`

	// When
	var rule CloseRule
	err := json.Unmarshal([]byte(input), &rule)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Selector != ".btn-secondary" {
		t.Errorf("expected selector '.btn-secondary', got '%s'", rule.Selector)
	}
	if rule.Match != "text" {
		t.Errorf("expected match 'text', got '%s'", rule.Match)
	}
	if rule.Value != "Close" {
		t.Errorf("expected value 'Close', got '%s'", rule.Value)
	}
}

func TestLoadConfig_FileNotExist(t *testing.T) {
	// Given
	configPath := "/nonexistent/path/connect.json"

	// When
	cfg, err := LoadConfig(configPath)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Links) != 0 {
		t.Errorf("expected empty links, got %d", len(cfg.Links))
	}
	if len(cfg.Modals) != 0 {
		t.Errorf("expected empty modals, got %d", len(cfg.Modals))
	}
	if len(cfg.Close) != 0 {
		t.Errorf("expected empty close, got %d", len(cfg.Close))
	}
	if len(cfg.Ignore) != 0 {
		t.Errorf("expected empty ignore, got %d", len(cfg.Ignore))
	}
}

func TestLoadConfig_ValidConfig(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "connect.json")
	configContent := `{
		"links": [
			{
				"selector": ".nav-item",
				"match": "text",
				"mapping": {"Home": "home"}
			}
		],
		"modals": ["tags"],
		"close": [".modal-close"],
		"ignore": ["mockup"]
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// When
	cfg, err := LoadConfig(configPath)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Links) != 1 {
		t.Errorf("expected 1 link, got %d", len(cfg.Links))
	}
	if cfg.Links[0].Selector != ".nav-item" {
		t.Errorf("expected selector '.nav-item', got '%s'", cfg.Links[0].Selector)
	}
	if !reflect.DeepEqual(cfg.Modals, []string{"tags"}) {
		t.Errorf("expected modals ['tags'], got %v", cfg.Modals)
	}
	if len(cfg.Close) != 1 {
		t.Errorf("expected 1 close rule, got %d", len(cfg.Close))
	}
	if cfg.Close[0].Selector != ".modal-close" {
		t.Errorf("expected close selector '.modal-close', got '%s'", cfg.Close[0].Selector)
	}
	if !reflect.DeepEqual(cfg.Ignore, []string{"mockup"}) {
		t.Errorf("expected ignore ['mockup'], got %v", cfg.Ignore)
	}
}

func TestLoadConfig_MixedCloseRules(t *testing.T) {
	// Given
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "connect.json")
	configContent := `{
		"close": [
			".modal-close",
			{"selector": ".btn", "match": "text", "value": "Cancel"}
		]
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// When
	cfg, err := LoadConfig(configPath)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Close) != 2 {
		t.Fatalf("expected 2 close rules, got %d", len(cfg.Close))
	}
	if cfg.Close[0].Selector != ".modal-close" {
		t.Errorf("expected first close selector '.modal-close', got '%s'", cfg.Close[0].Selector)
	}
	if cfg.Close[0].Match != "" {
		t.Errorf("expected first close match '', got '%s'", cfg.Close[0].Match)
	}
	if cfg.Close[1].Selector != ".btn" {
		t.Errorf("expected second close selector '.btn', got '%s'", cfg.Close[1].Selector)
	}
	if cfg.Close[1].Match != "text" {
		t.Errorf("expected second close match 'text', got '%s'", cfg.Close[1].Match)
	}
	if cfg.Close[1].Value != "Cancel" {
		t.Errorf("expected second close value 'Cancel', got '%s'", cfg.Close[1].Value)
	}
}

func TestConfigPath(t *testing.T) {
	// Given
	designsDir := "/path/to/designs"

	// When
	result := ConfigPath(designsDir)

	// Then
	expected := "/path/to/designs/connect.json"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}
