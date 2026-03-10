package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {
	BusName    string // message bus name
	SourceName string // message source name
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.BusName, err = ensureSetAndNonEmpty("MESSAGE_BUS")
	if err != nil {
		return nil, err
	}

	cfg.SourceName, err = ensureSetAndNonEmpty("MESSAGE_SOURCE")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[CONFIG] BusName    = [%s]\n", cfg.BusName)
	fmt.Printf("[CONFIG] SourceName = [%s]\n", cfg.SourceName)

	return &cfg, nil
}

//
// end of file
//
