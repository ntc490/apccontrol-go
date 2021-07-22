package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strconv"
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
	LastPort    int     `yaml:"last_port"`
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

func (config *ConfigFile) AliasNameByNum(num int) (alias string, err error) {
	for _, alias := range config.Aliases {
		if alias.Port == num {
			return alias.Name, nil
		}
	}
	return "", errors.New("Invalid alias port num")
}

func (config *ConfigFile) AliasNumByName(name string) (port int, err error) {
	for _, alias := range config.Aliases {
		if alias.Name == name {
			return alias.Port, nil
		}
	}
	return 0, errors.New("Invalid alias port name")
}

// Might make sense to return a different value if we add a new alias
// or change an existing one
func (config *ConfigFile) SetAlias(num int, name string) (err error) {
	aliasNum, err := config.AliasNumByName(name)
	if err == nil && num != aliasNum {
		return errors.New("Alias name already in use for port " + strconv.Itoa(aliasNum))
	}
	for index, alias := range config.Aliases {
		if alias.Port == num {
			alias.Name = name
			alias.Description = ""
			config.Aliases[index] = alias
			return nil
		}
	}
	alias := Alias{num, name, ""}
	config.Aliases = append(config.Aliases, alias)
	return nil
}

func (config *ConfigFile) RmAliasByName(name string) (err error) {
	for index, alias := range config.Aliases {
		if alias.Name == name {
			config.rmAliasIndex(index)
			return nil
		}
	}
	return errors.New("Invalid alias port name")
}

func (config *ConfigFile) CheckBasicSettings() (err error) {
	if config.Hostname == "" {
		return errors.New("Config file must contain APC hostname")
	}
	if config.User == "" {
		return errors.New("Config file must contain APC device user")
	}
	return nil
}

func (config *ConfigFile) rmAliasIndex(index int) {
	config.Aliases = append(config.Aliases[:index], config.Aliases[index+1:]...)
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
