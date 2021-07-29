package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"log"
	"github.com/ziutek/telnet"
)

const OUTLET_MANAGER = "1"
const CONTROL_PORT = "1"
const PORT_ON_COMMAND = "1"
const PORT_OFF_COMMAND = "2"
const PORT_RESET_COMMAND = "3"

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
	apc.config.SetLastPort(num)
	err = apc.controlOutlet(num, PORT_ON_COMMAND)
	return err
}

func (apc *ApcConnector) Off(port string) (err error) {
	num, err := apc.portNumFromString(port)
	if err != nil {
		return err
	}
	fmt.Println("Turning off port:", num)
	apc.config.SetLastPort(num)
	err = apc.controlOutlet(num, PORT_OFF_COMMAND)
	return err
}

func (apc *ApcConnector) Reset(port string) (err error) {
	num, err := apc.portNumFromString(port)
	if err != nil {
		return err
	}
	fmt.Println("Reset port:", num)
	apc.config.SetLastPort(num)
	err = apc.controlOutlet(num, PORT_RESET_COMMAND)
	return err
}

func (apc *ApcConnector) portNumFromString(port string) (num int, err error) {
	// An empty string should defer to the LastPort var
	if port == "" {
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

func (apc *ApcConnector) controlOutlet(port int, command string) (err error) {
	server := apc.config.Hostname + ":23"

	t, err := telnet.Dial("tcp", server)
	if err != nil {
		//log.Fatalln("Error:", err)
		return err
	}

	t.SetUnixWriteMode(true)
	expect(t, "User Name :")
	sendln(t, apc.config.User)
	expect(t, "Password  :")
	sendln(t, apc.config.Password)

	expect(t, "<ESC>")
	expect(t, ">")
	sendln(t, OUTLET_MANAGER)

	expect(t, "<ESC>")
	expect(t, ">")
	sendln(t, strconv.Itoa(port))

	expect(t, "<ESC>")
	expect(t, ">")
	sendln(t, CONTROL_PORT)

	expect(t, "<ESC>")
	expect(t, ">")
	sendln(t, command)

	expect(t, "cancel :")
	sendln(t, "YES")

	return nil
}

const timeout = 10 * time.Second

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func expect(t *telnet.Conn, d ...string) {
	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}

func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}
