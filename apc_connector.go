package main

import (
	"errors"
	"fmt"
	"strconv"
)

type ApcConnector struct {
	config *ConfigFile
}

func NewApcConnectionFromConfigFile(config *ConfigFile) (apc *ApcConnector) {
	apc = &ApcConnector{}
	apc.config = config
	return apc
}

func (apc *ApcConnector) On(port string) (err error) {
	num, err := apc.portNumFromString(port)
	if err != nil {
		return err
	}
	fmt.Println("Turning on port:", num)
	// If successful - update the last port
	apc.config.LastPort = num
	return nil
}

func (apc *ApcConnector) Off(port string) (err error) {
	num, err := apc.portNumFromString(port)
	if err != nil {
		return err
	}
	fmt.Println("Turning off port:", num)
	// If successful - update the last port
	apc.config.LastPort = num
	return nil
}

func (apc *ApcConnector) Reset(port string) (err error) {
	num, err := apc.portNumFromString(port)
	if err != nil {
		return err
	}
	fmt.Println("Reset port:", num)
	// If successful - update the last port
	apc.config.LastPort = num
	return nil
}

func (apc *ApcConnector) portNumFromString(port string) (num int, err error) {
	// An empty string should defer to the LastPort var
	if port == "" {
		fmt.Println("Last port is", apc.config.LastPort)
		return apc.config.LastPort, nil
	}

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
