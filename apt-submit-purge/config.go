package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {
	// asset cache details
	AssetBucket     string
	AssetFilesystem string

	// event bus definitions
	//BusName        string // the event bus name
	//BusEventSource string // the source of published events

}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config
	var err error

	// asset cache details
	cfg.AssetBucket, err = ensureSetAndNonEmpty("ASSET_BUCKET")
	if err != nil {
		return nil, err
	}
	cfg.AssetFilesystem, err = ensureSetAndNonEmpty("ASSET_FILESYSTEM")
	if err != nil {
		return nil, err
	}

	// event bus definitions
	//cfg.BusName = envWithDefault("EVENT_BUS_NAME", "")
	//cfg.BusEventSource = envWithDefault("EVENT_SRC_NAME", "")

	// asset cache details
	fmt.Printf("[CONFIG] AssetBucket     = [%s]\n", cfg.AssetBucket)
	fmt.Printf("[CONFIG] AssetFilesystem = [%s]\n", cfg.AssetFilesystem)

	// event bus definitions
	//fmt.Printf("[CONFIG] BusName         = [%s]\n", cfg.BusName)
	//fmt.Printf("[CONFIG] BusEventSource  = [%s]\n", cfg.BusEventSource)

	return &cfg, nil
}

//
// end of file
//
