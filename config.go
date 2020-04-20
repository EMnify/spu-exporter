package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	Prometheus struct {
		Outfile string `yaml:"outfile"`
	} `yaml:"prometheus"`
	Ssh struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		User    string `yaml:"user"`
		Keyfile string `yaml:keyfile`
		Command string `yaml:command`
	} `yaml:"ssh"`
}

func readConfig(filename string) AppConfig {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var cfg AppConfig
	err = cfg.Parse(yamlFile)

	if err != nil {
		log.Fatalf("Failed to unmarshal: %v", err)
	}
	fmt.Printf("%+v\n", cfg)
	return cfg
}

func (c *AppConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	// Optional add checks for required fields

	return nil
}
