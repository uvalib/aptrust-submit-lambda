//
// main message processing
//

package main

import (
	"encoding/json"
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

	// get the client details
	_, err = dao.GetBagsByStatus(BagStatusPendingIngest)
	if err != nil {
		return err
	}

	// do more stuff

	return nil
}

//
// end of file
//
