package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Service struct {
	Port int
	Text map[string]string
}

type Node struct {
	Host string
	Interfaces []string
	Services map[string]Service
}

type Config struct {
	Nodes map[string]Node
}

func (c *Config) Load(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}

	return nil
}
