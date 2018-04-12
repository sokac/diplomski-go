package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"syscall"
)

const TEMPLATE string = `# Managed by Manager
upstream manager_backend {
    server 127.0.0.1:%d;
}
`

type NginxConfigurator struct {
	configPath string
	pidFile    string
}

var _ Subscribers = &NginxConfigurator{} // enforce type

func NewNginxConfiguration(pidFile, configPath string) *NginxConfigurator {
	return &NginxConfigurator{
		configPath: configPath,
		pidFile:    pidFile,
	}
}

func (c *NginxConfigurator) readPID() (int, error) {
	o, err := ioutil.ReadFile(c.pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(o)))
}

func (c *NginxConfigurator) reload() error {
	pid, err := c.readPID()
	if err != nil {
		return err
	}
	return syscall.Kill(pid, syscall.SIGHUP)
}

func (c *NginxConfigurator) NewContainer(port int) {
	content := fmt.Sprintf(TEMPLATE, port)
	err := ioutil.WriteFile(c.configPath, []byte(content), 0644)
	if err != nil {
		log.Println("Error writing to nginx file:", err.Error())
	}
	err = c.reload()
	if err != nil {
		log.Println("Error sending a SIGHUP to nginx:", err.Error())
	}
}
