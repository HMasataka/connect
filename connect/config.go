package connect

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Links  []LinkRule  `json:"links"`
	Modals []string    `json:"modals"`
	Close  []CloseRule `json:"close"`
	Ignore []string    `json:"ignore"`
}

type LinkRule struct {
	Selector string            `json:"selector"`
	Match    string            `json:"match"`
	Mapping  map[string]string `json:"mapping"`
	Target   string            `json:"target"`
}

type CloseRule struct {
	Selector string `json:"selector"`
	Match    string `json:"match"`
	Value    string `json:"value"`
}

func (c *CloseRule) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		c.Selector = str
		c.Match = ""
		c.Value = ""
		return nil
	}

	type closeRuleAlias CloseRule
	var alias closeRuleAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*c = CloseRule(alias)
	return nil
}

func ConfigPath(designsDir string) string {
	return filepath.Join(designsDir, "connect.json")
}

func LoadConfig(configPath string) (Config, error) {
	cfg := Config{
		Links:  []LinkRule{},
		Modals: []string{},
		Close:  []CloseRule{},
		Ignore: []string{},
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	if cfg.Links == nil {
		cfg.Links = []LinkRule{}
	}
	if cfg.Modals == nil {
		cfg.Modals = []string{}
	}
	if cfg.Close == nil {
		cfg.Close = []CloseRule{}
	}
	if cfg.Ignore == nil {
		cfg.Ignore = []string{}
	}

	return cfg, nil
}
