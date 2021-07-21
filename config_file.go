package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type Alias struct {
	Port        int    `yaml:"port"`
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type ConfigFile struct {
	Filename    string
	Hostname    string  `yaml:"hostname"`
	User        string  `yaml:"user"`
	Password    string  `yaml:"password,omitempty"`
	LastPort    string  `yaml:"last_port"`
	Description string  `yaml:"description,omitempty"`
	Aliases     []Alias `yaml:"aliases"`
}

func NewConfigFile(filename string) (config *ConfigFile) {
	config = &ConfigFile{}
	config.Filename = filename
	return config
}

func (config *ConfigFile) Read() (err error) {
	filename, err := expandUserDir(config.Filename)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigFile) Write() (err error) {
	config.Description = "modified"
	data, err := yaml.Marshal(config)
	filename, err := expandUserDir(config.Filename)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigFile) AliasByNum(num int) (alias string, err error) {
	for _, alias := range config.Aliases {
		if alias.Port == num {
			return alias.Name, nil
		}
	}
	return "", errors.New("Invalid alias port num")
}

func (config *ConfigFile) AliasByName(name string) (port int, err error) {
	for _, alias := range config.Aliases {
		if alias.Name == name {
			return alias.Port, nil
		}
	}
	return 0, errors.New("Invalid alias port name")
}

// Expand ~ to $HOME in config file spec
func expandUserDir(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
