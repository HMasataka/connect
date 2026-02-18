package connect

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Selectors Selectors                    `json:"selectors"`
	Mapping   map[string]string            `json:"mapping"`
	Toolbar   map[string]string            `json:"toolbar"`
	Custom    map[string]map[string]string `json:"custom"`
	Ignore    []string                     `json:"ignore"`
}

type Selectors struct {
	Nav            string `json:"nav"`
	Toolbar        string `json:"toolbar"`
	ActiveClass    string `json:"activeClass"`
	ModalClose     string `json:"modalClose"`
	ModalCloseText string `json:"modalCloseText"`
}

func DefaultSelectors() Selectors {
	return Selectors{
		Nav:            ".nav-item",
		Toolbar:        ".toolbar-btn",
		ActiveClass:    "active",
		ModalClose:     ".modal-close",
		ModalCloseText: ".modal-footer .btn-secondary",
	}
}

func DefaultConfig() Config {
	return Config{
		Selectors: DefaultSelectors(),
		Mapping:   make(map[string]string),
		Toolbar:   make(map[string]string),
		Custom:    make(map[string]map[string]string),
		Ignore:    []string{},
	}
}

func ConfigPath(designsDir string) string {
	return filepath.Join(designsDir, "connect.json")
}

func LoadConfig(configPath string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return cfg, err
	}

	if fileCfg.Selectors.Nav != "" {
		cfg.Selectors.Nav = fileCfg.Selectors.Nav
	}
	if fileCfg.Selectors.Toolbar != "" {
		cfg.Selectors.Toolbar = fileCfg.Selectors.Toolbar
	}
	if fileCfg.Selectors.ActiveClass != "" {
		cfg.Selectors.ActiveClass = fileCfg.Selectors.ActiveClass
	}
	if fileCfg.Selectors.ModalClose != "" {
		cfg.Selectors.ModalClose = fileCfg.Selectors.ModalClose
	}
	if fileCfg.Selectors.ModalCloseText != "" {
		cfg.Selectors.ModalCloseText = fileCfg.Selectors.ModalCloseText
	}
	if fileCfg.Mapping != nil {
		cfg.Mapping = fileCfg.Mapping
	}
	if fileCfg.Toolbar != nil {
		cfg.Toolbar = fileCfg.Toolbar
	}
	if fileCfg.Custom != nil {
		cfg.Custom = fileCfg.Custom
	}
	if fileCfg.Ignore != nil {
		cfg.Ignore = fileCfg.Ignore
	}

	return cfg, nil
}
