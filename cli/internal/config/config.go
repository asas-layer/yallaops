package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the CLI's local configuration, loaded from ~/.yallaops/config.yaml.
type Config struct {
	CurrentContext string             `yaml:"current_context"`
	Contexts       map[string]Context `yaml:"contexts"`
}

// Context is a named control-plane endpoint plus the credentials and default
// environment to use when talking to it, analogous to a kubeconfig context.
type Context struct {
	Endpoint       string `yaml:"endpoint"`
	Token          string `yaml:"token,omitempty"`
	DefaultService string `yaml:"default_service,omitempty"`
}

// Default returns the fallback config used when no config file exists yet.
func Default() *Config {
	return &Config{
		CurrentContext: "default",
		Contexts: map[string]Context{
			"default": {Endpoint: "localhost:50051"},
		},
	}
}

// Path returns the path to the CLI config file, honoring $YALLAOPS_CONFIG.
func Path() (string, error) {
	if p := os.Getenv("YALLAOPS_CONFIG"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config: resolve home dir: %w", err)
	}
	return filepath.Join(home, ".yallaops", "config.yaml"), nil
}

// Load reads the config file, returning Default() if it does not exist.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Default(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}
	return &cfg, nil
}

// Save writes the config file, creating ~/.yallaops if needed.
func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("config: create dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("config: write %s: %w", path, err)
	}
	return nil
}

// Current returns the active context, erroring if it isn't defined.
func (c *Config) Current() (Context, error) {
	ctx, ok := c.Contexts[c.CurrentContext]
	if !ok {
		return Context{}, fmt.Errorf("config: context %q not found", c.CurrentContext)
	}
	return ctx, nil
}
