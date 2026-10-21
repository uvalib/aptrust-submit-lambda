package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {

	// general definitions
	ReminderDays int // number of days after the original email to send a reminder

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

// these settings are meaningless at or below zero and misbehave in ways that are hard
// to spot at runtime, so reject them up front
func ensurePositive(env string, value int) error {

	if value < 1 {
		err := fmt.Errorf("environment variable must be a positive integer: [%s] = [%d]", env, value)
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	return nil
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	// general definitions
	cfg.ReminderDays, err = envToInt("REMINDER_DAYS")
	if err != nil {
		return nil, err
	}

	err = ensurePositive("REMINDER_DAYS", cfg.ReminderDays)
	if err != nil {
		return nil, err
	}

	// event bus definitions
	cfg.BusName, err = ensureSetAndNonEmpty("EVENT_BUS_NAME")
	if err != nil {
		return nil, err
	}
	cfg.BusEventSource, err = ensureSetAndNonEmpty("EVENT_SRC_NAME")
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

	// general definitions
	fmt.Printf("[CONFIG] ReminderDays    = [%d]\n", cfg.ReminderDays)

	// event bus definitions
	fmt.Printf("[CONFIG] BusName         = [%s]\n", cfg.BusName)
	fmt.Printf("[CONFIG] BusEventSource  = [%s]\n", cfg.BusEventSource)

	// database configuration
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
