//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	fmt.Printf("INFO: event %s from %s -> %s\n", messageId, messageSrc, string(rawMsg))

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	busCfg := uvaaptsbus.UvaBusConfig{
		Source:  cfg.BusEventSource,
		BusName: cfg.BusName,
		Log:     log.Default(),
	}

	// create message bus client
	bus, err := uvaaptsbus.NewUvaBus(busCfg)
	if err != nil {
		fmt.Printf("ERROR: creating event bus client (%s)\n", err.Error())
		return err
	}
	fmt.Printf("Using: %s@%s\n", cfg.BusEventSource, cfg.BusName)

	// create event
	ev := uvaaptsbus.UvaBusEvent{}
	ev.EventName = uvaaptsbus.EventScheduleCheckPending
	//ev.ClientId = "none"
	//ev.SubmissionId = "none"
	//ev.BagId = "none"

	// publish event
	err = bus.PublishEvent(&ev)
	if err != nil {
		fmt.Printf("ERROR: publishing event (%s)\n", err.Error())
		return err
	}

	return nil
}

//
// end of file
//
