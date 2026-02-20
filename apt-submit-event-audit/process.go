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

	// convert to uvaaptsbus event
	ev, err := uvaaptsbus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: un-marshaling bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: EVENT %s from %s -> %s\n", messageId, messageSrc, ev.String())
	return nil
}

//
// end of file
//
