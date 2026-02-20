//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/apts-bus-definitions/uvaaptsbus"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	fmt.Printf("INFO: EVENT %s from %s -> %s\n", messageId, messageSrc, string(rawMsg))

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	busCfg := uvaaptsbus.UvaBusConfig{
		Source:  cfg.SourceName,
		BusName: cfg.BusName,
	}

	// create message bus client
	bus, err := uvaaptsbus.NewUvaBus(busCfg)
	if err != nil {
		fmt.Printf("ERROR: creating event bus client (%s)\n", err.Error())
		return err
	}
	fmt.Printf("Using: %s@%s\n", cfg.SourceName, cfg.BusName)

	// create event
	ev := uvaaptsbus.UvaBusEvent{}
	ev.EventName = uvaaptsbus.EventScheduleCheckPending
	ev.ClientId = "none"
	ev.SubmissionId = "none"
	ev.BagId = "none"

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
