//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// create the data access object
	dao, err := newDao(cfg)
	if err != nil {
		return err
	}

	// cleanup on exit
	defer dao.Close()

	// get the bags
	bags, err := dao.GetBagsByStatus(BagStatusPendingIngest)
	if err != nil {
		return err
	}

	if len(bags) != 0 {
		for _, b := range bags {
			fmt.Printf("DEBUG: checking APT for ingest status of '%s' (%s))\n", b.Name, b.Identifier)
		}
	} else {
		fmt.Printf("INFO: no bags in '%s' status)\n", BagStatusPendingIngest)
	}

	return nil
}

//
// end of file
//
