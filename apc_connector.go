package main

import (
	"errors"
	"fmt"
	"strconv"
)

type ApcConnector struct {
	Host     string
	User     string
	Password string
	aliases  []Alias
	LastPort string
	config   *ConfigFile
}

func NewApcConnectionFromConfigFile(config *ConfigFile) (apc *ApcConnector) {
	apc = &ApcConnector{}
	apc.Host = config.Hostname
	apc.User = config.User
	apc.Password = config.Password
	apc.LastPort = config.LastPort
	apc.config = config
	return apc
}

func (apc *ApcConnector) On(port string) (err error) {
	num, err := apc.portNumFromString(port)
	if err != nil {
		return err
	}
	fmt.Println("Turning on port:", num)
	return nil
}

func (apc *ApcConnector) Off(port string) (err error) {
	return nil
}

func (apc *ApcConnector) Reset(port string) (err error) {
	return nil
}

func (apc *ApcConnector) portNumFromString(port string) (num int, err error) {
	// Check if the user simply passed in a port number
	if num, err := strconv.Atoi(port); err == nil {
		return num, nil
	}

	// Can we convert the alias name string to a port number?
	if num, err = apc.config.AliasNumByName(port); err == nil {
		return num, nil
	}

	// Couldn't do anything with the port string
	return 0, errors.New("Unable to decode port '" + port + "'")
}
