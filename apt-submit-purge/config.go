package main

import (
	"fmt"
)

// Config defines all the service configuration parameters
type Config struct {
	// asset cache details
	AssetBucket     string
	AssetFilesystem string
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

	// asset cache details
	fmt.Printf("[CONFIG] AssetBucket     = [%s]\n", cfg.AssetBucket)
	fmt.Printf("[CONFIG] AssetFilesystem = [%s]\n", cfg.AssetFilesystem)

	return &cfg, nil
}

//
// end of file
//
