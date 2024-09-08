package global

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	GlobalConfigDirectory = ".justjump"
	GlobalConfigFile      = "registry.yaml"
)

type Jumproot struct {
	Name string `yaml:"-"`
	Root string `yaml:"jumproot"`
}

type JumpRoots map[string]Jumproot

type GlobalConfig interface {
	Load() error
	Save() error
	JumpRoots() JumpRoots
	ObtainJumpRoot(name string) (Jumproot, bool)
	RegisterJumpRoot(name string, project Jumproot) error
	DeleteJumpRoot(name string) error
}

func IsGlobalConfigPresent() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configPath := filepath.Join(home, GlobalConfigDirectory, GlobalConfigFile)
	_, err = os.Stat(configPath)
	return !os.IsNotExist(err)
}

type globalConfig struct {
	jumpRoots JumpRoots
}

func New() (GlobalConfig, error) {
	if !IsGlobalConfigPresent() {
		return &globalConfig{
			jumpRoots: make(JumpRoots),
		}, nil
	} else {
		// load the existing config
		config := &globalConfig{}
		err := config.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}

		return config, nil
	}
}

func (c *globalConfig) Load() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(home, GlobalConfigDirectory, GlobalConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &c.jumpRoots)
	if err != nil {
		return fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	return nil
}

func (c *globalConfig) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	var configFolder string = filepath.Join(home, GlobalConfigDirectory)
	err = os.MkdirAll(configFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create JustJump config directory: %w", err)
	}

	configPath := filepath.Join(home, GlobalConfigDirectory, GlobalConfigFile)
	data, err := yaml.Marshal(c.jumpRoots)
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *globalConfig) JumpRoots() JumpRoots {
	return c.jumpRoots
}

func (c *globalConfig) ObtainJumpRoot(name string) (Jumproot, bool) {
	project, exists := c.jumpRoots[name]
	return project, exists
}

func (c *globalConfig) RegisterJumpRoot(name string, project Jumproot) error {
	if _, exists := c.jumpRoots[name]; exists {
		return fmt.Errorf("project %s already exists", name)
	}

	c.jumpRoots[name] = project
	return nil
}

func (c *globalConfig) DeleteJumpRoot(name string) error {
	if _, exists := c.jumpRoots[name]; !exists {
		return fmt.Errorf("project %s does not exist", name)
	}

	delete(c.jumpRoots, name)
	return nil
}
