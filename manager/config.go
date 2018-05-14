package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	DockerImage            string
	NginxPIDFile           string
	NginxConfiguration     string
	VersionOverlapDuration int
}

func (c *Config) Validate() error {
	if c.DockerImage == "" {
		return fmt.Errorf("Docker image not defined")
	}
	if c.NginxConfiguration == "" {
		return fmt.Errorf("Nginx Configuration file not defined")
	}
	return nil
}

func loadConfig(f string) (*Config, error) {
	c := &Config{}
	fileData, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileData, c)
	if err != nil {
		return nil, err
	}
	if err := c.Validate(); c != nil {
		return nil, err
	}

	return c, err
}
