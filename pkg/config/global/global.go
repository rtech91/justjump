package global

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigDirectory = ".config/justjump"
	JumpRootsDir    = "jumproots.d"
	FileExtension   = ".jpath"
	LastJumpFile    = "last_jump"
)

type Jumproot struct {
	Name string
	Root string
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

	configPath := filepath.Join(home, ConfigDirectory, JumpRootsDir)
	_, err = os.Stat(configPath)
	return !errors.Is(err, fs.ErrNotExist)
}

type globalConfig struct {
	jumpRoots JumpRoots
}

func New() (GlobalConfig, error) {
	config := &globalConfig{
		jumpRoots: make(JumpRoots),
	}
	err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return config, nil
}

func (c *globalConfig) Load() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	jumpRootsPath := filepath.Join(home, ConfigDirectory, JumpRootsDir)
	err = os.MkdirAll(jumpRootsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	files, err := os.ReadDir(jumpRootsPath)
	if err != nil {
		return fmt.Errorf("failed to read jump roots directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), FileExtension) {
			name := strings.TrimSuffix(file.Name(), FileExtension)
			filePath := filepath.Join(jumpRootsPath, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read jump root file %s: %w", filePath, err)
			}
			absolutePath := strings.TrimSpace(string(data))
			if !filepath.IsAbs(absolutePath) {
				return fmt.Errorf("invalid path in file %s: must be an absolute path", filePath)
			}
			c.jumpRoots[name] = Jumproot{Name: name, Root: absolutePath}
		}
	}

	return nil
}

func (c *globalConfig) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	jumpRootsPath := filepath.Join(home, ConfigDirectory, JumpRootsDir)
	err = os.MkdirAll(jumpRootsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	for name, jumproot := range c.jumpRoots {
		filePath := filepath.Join(jumpRootsPath, name+FileExtension)
		err := os.WriteFile(filePath, []byte(jumproot.Root), 0644)
		if err != nil {
			return fmt.Errorf("failed to write jump root file %s: %w", filePath, err)
		}
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

	if !filepath.IsAbs(project.Root) {
		return fmt.Errorf("invalid path for project %s: must be an absolute path", name)
	}

	c.jumpRoots[name] = project

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(home, ConfigDirectory, JumpRootsDir, name+FileExtension)
	err = os.WriteFile(filePath, []byte(project.Root), 0644)
	if err != nil {
		return fmt.Errorf("failed to write jump root file %s: %w", filePath, err)
	}

	return nil
}

func (c *globalConfig) DeleteJumpRoot(name string) error {
	if _, exists := c.jumpRoots[name]; !exists {
		return fmt.Errorf("project %s does not exist", name)
	}

	delete(c.jumpRoots, name)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(home, ConfigDirectory, JumpRootsDir, name+FileExtension)
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete jump root file %s: %w", filePath, err)
	}

	return nil
}
