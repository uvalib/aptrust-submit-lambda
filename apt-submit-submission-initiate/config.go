package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {
	// ingest details
	InboundBucket string

	// event bus definitions
	BusName        string // the event bus name
	BusEventSource string // the source of published events

	// database details
	DbHost     string // database host
	DbPort     int    // database port
	DbName     string // database name
	DbUser     string // database user
	DbPassword string // database password
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config
	var err error

	// ingest details
	cfg.InboundBucket, err = ensureSetAndNonEmpty("INBOUND_BUCKET")
	if err != nil {
		return nil, err
	}

	// event bus definitions
	cfg.BusName = envWithDefault("EVENT_BUS_NAME", "")
	cfg.BusEventSource = envWithDefault("EVENT_SRC_NAME", "")

	// database details
	cfg.DbHost, err = ensureSetAndNonEmpty("DB_HOST")
	if err != nil {
		return nil, err
	}
	cfg.DbPort, err = envToInt("DB_PORT")
	if err != nil {
		return nil, err
	}
	cfg.DbName, err = ensureSetAndNonEmpty("DB_NAME")
	if err != nil {
		return nil, err
	}
	cfg.DbUser, err = ensureSetAndNonEmpty("DB_USER")
	if err != nil {
		return nil, err
	}
	cfg.DbPassword, err = ensureSetAndNonEmpty("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	// ingest details
	fmt.Printf("[CONFIG] InboundBucket  = [%s]\n", cfg.InboundBucket)

	// event bus definitions
	fmt.Printf("[CONFIG] BusName        = [%s]\n", cfg.BusName)
	fmt.Printf("[CONFIG] BusEventSource = [%s]\n", cfg.BusEventSource)

	// database details
	fmt.Printf("[CONFIG] DbHost         = [%s]\n", cfg.DbHost)
	fmt.Printf("[CONFIG] DbPort         = [%d]\n", cfg.DbPort)
	fmt.Printf("[CONFIG] DbName         = [%s]\n", cfg.DbName)
	fmt.Printf("[CONFIG] DbUser         = [%s]\n", cfg.DbUser)
	fmt.Printf("[CONFIG] DbPassword     = [REDACTED]\n")

	return &cfg, nil
}

//
// end of file
//
