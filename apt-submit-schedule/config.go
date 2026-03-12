package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {
	BusName        string // message bus name
	BusEventSource string // message source name
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.BusName, err = ensureSetAndNonEmpty("EVENT_BUS_NAME")
	if err != nil {
		return nil, err
	}

	cfg.BusEventSource, err = ensureSetAndNonEmpty("EVENT_SRC_NAME")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[CONFIG] BusName        = [%s]\n", cfg.BusName)
	fmt.Printf("[CONFIG] BusEventSource = [%s]\n", cfg.BusEventSource)

	return &cfg, nil
}

//
// end of file
//
