package local

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	LocalConfigFile = ".justjump.yaml"
)

var (
	ErrRelPathNotExist = errors.New("relative path does not exist")
)

type LocalConfig interface {
	Load() error
	Save() error
	JumpPoints() []string
	AddJumpPoint(path string) error
	RemoveJumpPoint(path string) error
}

func IsLocalConfigPresent(jumproot string) bool {
	localConfigPath := jumproot + "/" + LocalConfigFile
	if _, err := os.Stat(localConfigPath); errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}

type JumpPointsList struct {
	Paths []string `yaml:"jumppoints"`
}

type localConfig struct {
	jumpRoot        string
	localConfigPath string
	jumpPoints      JumpPointsList
}

func New(jumproot string) (LocalConfig, error) {

	if !IsLocalConfigPresent(jumproot) {
		return &localConfig{
			jumpRoot:        jumproot,
			localConfigPath: jumproot + "/" + LocalConfigFile,
			jumpPoints: JumpPointsList{
				Paths: make([]string, 0),
			},
		}, nil
	}

	// load the existing config
	config := &localConfig{
		jumpRoot:        jumproot,
		localConfigPath: jumproot + "/" + LocalConfigFile,
		jumpPoints: JumpPointsList{
			Paths: make([]string, 0),
		},
	}
	err := config.Load()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *localConfig) Load() error {
	data, err := os.ReadFile(c.localConfigPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &c.jumpPoints)
	if err != nil {
		return err
	}

	return nil
}

func (c *localConfig) Save() error {
	data, err := yaml.Marshal(c.jumpPoints)
	if err != nil {
		return err
	}

	err = os.WriteFile(c.localConfigPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *localConfig) JumpPoints() []string {
	return c.jumpPoints.Paths
}

func (c *localConfig) AddJumpPoint(rpath string) error {
	// Ensure the path ends with a slash
	if rpath[len(rpath)-1] != '/' {
		rpath += "/"
	}

	// Check if the relative path exists within the jump root
	if err := c.CheckLocalPathExist(rpath); err != nil {
		return err
	}

	// Avoid adding the relative path if it already exists in the configuration
	for _, existingPath := range c.jumpPoints.Paths {
		if existingPath == rpath {
			return errors.New("jump point already exists")
		}
	}

	// Add the relative path to the list
	c.jumpPoints.Paths = append(c.jumpPoints.Paths, rpath)

	// Save the updated configuration
	err := c.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *localConfig) RemoveJumpPoint(path string) error {
	// Normalize to absolute path
	if !filepath.IsAbs(path) {
		path = filepath.Join(c.jumpRoot, path)
	}
	path = filepath.Clean(path) + "/"

	bestIndex := -1
	bestMatchLen := -1

	for i, relPath := range c.jumpPoints.Paths {
		absJumpPoint := filepath.Clean(filepath.Join(c.jumpRoot, relPath)) + "/"

		if strings.HasPrefix(path, absJumpPoint) {
			if len(absJumpPoint) > bestMatchLen {
				bestMatchLen = len(absJumpPoint)
				bestIndex = i
			}
		}
	}

	if bestIndex != -1 {
		c.jumpPoints.Paths = append(c.jumpPoints.Paths[:bestIndex], c.jumpPoints.Paths[bestIndex+1:]...)

		err := c.Save()
		if err != nil {
			return err
		}

		return nil
	}

	return ErrRelPathNotExist
}

func (c *localConfig) CheckLocalPathExist(rpath string) error {
	rpathstat, err := os.Stat(c.jumpRoot + "/" + rpath)
	if err != nil || !rpathstat.IsDir() {
		return ErrRelPathNotExist
	}

	return nil
}
