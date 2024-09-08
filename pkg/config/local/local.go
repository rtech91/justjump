package local

import (
	"errors"
	"log"
	"os"

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
	if _, err := os.Stat(localConfigPath); os.IsNotExist(err) {
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
	log.Println(c.jumpPoints)
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
	// check that path ends with a slash or is a directory
	if rpath[len(rpath)-1] != '/' {
		rpath += "/"
		if err := c.CheckLocalPathExist(rpath); err != nil {
			return err
		}
	}

	c.jumpPoints.Paths = append(c.jumpPoints.Paths, rpath)
	return nil
}

func (c *localConfig) RemoveJumpPoint(path string) error {
	for i, p := range c.jumpPoints.Paths {
		if p == path {
			c.jumpPoints.Paths = append(c.jumpPoints.Paths[:i], c.jumpPoints.Paths[i+1:]...)
			return nil
		}
	}

	return ErrRelPathNotExist
}

func (c *localConfig) CheckLocalPathExist(rpath string) error {
	rpathstat, err := os.Stat(c.jumpRoot + rpath)
	if os.IsNotExist(err) || !rpathstat.IsDir() {
		return ErrRelPathNotExist
	}

	return nil
}
