package config

import (
	"log"
	"os"
	"strconv"
)

type CacheConfig struct {
	Driver string
	Name   string
	Host   string
	Port   int
}

func (c *CacheConfig) Setup() {
	c.Driver = os.Getenv("CACHE_DRIVER")
	if c.Driver == "" {
		c.Driver = "default_driver" // Add a default value if needed
	}

	c.Name = os.Getenv("CACHE_NAME")
	if c.Name == "" {
		c.Name = "default_name"
	}

	c.Host = os.Getenv("CACHE_HOST")
	if c.Host == "" {
		c.Host = "localhost"
	}
	portStr := os.Getenv("CACHE_PORT")
	if portStr == "" {
		c.Port = 6379
	} else {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
		}
		c.Port = port
	}
}
