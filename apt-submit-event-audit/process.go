//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to uvaaptsbus event
	ev, err := uvaaptsbus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling bus event (%s)\n", err.Error())
		return nil
	}

	// make the workflow event
	if len(ev.Detail) != 0 {
		wf, err := uvaaptsbus.MakeWorkflowEvent(ev.Detail)
		if err != nil {
			fmt.Printf("ERROR: unmarshaling workflow event (%s)\n", err.Error())
			fmt.Printf("INFO: event %s from %s -> %s (%s)\n", messageId, messageSrc, ev.String(), ev.Detail)
		} else {
			fmt.Printf("INFO: event %s from %s -> %s/%s\n", messageId, messageSrc, ev.String(), wf.String())
		}
	} else {
		fmt.Printf("INFO: event %s from %s -> %s\n", messageId, messageSrc, ev.String())
	}

	return nil
}

//
// end of file
//
