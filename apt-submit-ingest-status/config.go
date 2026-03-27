package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {
	// general definitions
	MaxRequests int // maximum number of requests in any lambda iteration
	HttpTimeout int // http connection timeout

	// event bus definitions
	BusName        string // the event bus name
	BusEventSource string // the source of published events

	// APTrust configuration
	APTUser      string // the APTrust user name
	APTKey       string // the APTrust access key
	APTStatusUrl string // the APTrust status URL

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
	// general definitions
	cfg.MaxRequests, err = envToInt("MAX_REQUESTS")
	if err != nil {
		return nil, err
	}
	cfg.HttpTimeout, err = envToInt("HTTP_TIMEOUT")
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

	// APTrust definitions
	cfg.APTUser, err = ensureSetAndNonEmpty("APTRUST_USER")
	if err != nil {
		return nil, err
	}
	cfg.APTKey, err = ensureSetAndNonEmpty("APTRUST_KEY")
	if err != nil {
		return nil, err
	}
	cfg.APTStatusUrl, err = ensureSetAndNonEmpty("APTRUST_STATUS_URL")
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
	fmt.Printf("[CONFIG] MaxRequests     = [%d]\n", cfg.MaxRequests)
	fmt.Printf("[CONFIG] HttpTimeout     = [%d]\n", cfg.HttpTimeout)

	// event bus definitions
	fmt.Printf("[CONFIG] BusName         = [%s]\n", cfg.BusName)
	fmt.Printf("[CONFIG] BusEventSource  = [%s]\n", cfg.BusEventSource)

	// APTrust configuration
	fmt.Printf("[CONFIG] APTUser         = [%s]\n", cfg.APTUser)
	fmt.Printf("[CONFIG] APTKey          = [REDACTED]\n")
	fmt.Printf("[CONFIG] APTStatusUrl    = [%s]\n", cfg.APTStatusUrl)
	fmt.Printf("[CONFIG] HttpTimeout     = [%d]\n", cfg.HttpTimeout)

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
