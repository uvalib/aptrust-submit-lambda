package main

import (
	"fmt"
)

type Config struct {
	// event bus definitions
	BusName        string // the event bus name
	BusEventSource string // the source of published events

	// database configuration
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

	// event bus definitions
	cfg.BusName, err = ensureSetAndNonEmpty("MESSAGE_BUS")
	if err != nil {
		return nil, err
	}
	cfg.BusEventSource, err = ensureSetAndNonEmpty("MESSAGE_SOURCE")
	if err != nil {
		return nil, err
	}

	// database definitions
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

	// event bus definitions
	fmt.Printf("[CONFIG] BusName         = [%s]\n", cfg.BusName)
	fmt.Printf("[CONFIG] BusEventSource  = [%s]\n", cfg.BusEventSource)

	// database definitions
	fmt.Printf("[CONFIG] DbHost          = [%s]\n", cfg.DbHost)
	fmt.Printf("[CONFIG] DbPort          = [%d]\n", cfg.DbPort)
	fmt.Printf("[CONFIG] DbName          = [%s]\n", cfg.DbName)
	fmt.Printf("[CONFIG] DbUser          = [%s]\n", cfg.DbUser)
	fmt.Printf("[CONFIG] DbPassword      = [REDACTED]\n")

	return &cfg, nil
}

//
// end of file
//
