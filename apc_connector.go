package main

import (
//	"errors"
)

type ApcConnector struct {
	Host     string
	User     string
	Password string
	aliases  []Alias
	LastPort string
}

func ApcConnectorFromConfig(config *ConfigFile) (apc *ApcConnector) {
	apc = &ApcConnector{}
	apc.Host = config.Hostname
	apc.User = config.User
	apc.Password = config.Password
	apc.LastPort = config.LastPort
	return apc
}

func (apc *ApcConnector) ApcPortName(num int) string {
	for _, alias := range apc.aliases {
		if alias.Port == num {
			return alias.Name
		}
	}
	return "Unknown"
}

func (apc *ApcConnector) ApcPortNum(name string) int {
	for _, alias := range apc.aliases {
		if alias.Name == name {
			return alias.Port
		}
	}
	return -1
}

func (apc *ApcConnector) ApcOn() (err error) {
	return nil
}

func (apc *ApcConnector) ApcOff() (err error) {
	return nil
}

func (apc *ApcConnector) ApcReset() (err error) {
	return nil
}
