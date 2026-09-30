package main

import (
	"encoding/json"
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

func NewEventBus(eventBus string, eventSource string) (uvaaptsbus.UvaBus, error) {
	// we will accept bad config and return nil quietly
	if len(eventBus) == 0 {
		fmt.Printf("INFO: Event bus is not configured, no events emitted\n")
		return nil, uvaaptsbus.ErrConfig
	}

	cfg := uvaaptsbus.UvaBusConfig{BusName: eventBus, Source: eventSource, Log: nil}
	return uvaaptsbus.NewUvaBus(cfg)
}

func publishWorkflowEvent(bus uvaaptsbus.UvaBus, eventName string, clientId string, submissionId string, bagId string, extra string) error {

	// an absent bus is a deliberate configuration for running locally, so skipping the
	// publish is not a failure and must not be reported as one; callers need a non-nil
	// error to mean the event genuinely did not get out. Say something every time
	// though, because a deployed function reaching here stalls the workflow silently
	if bus == nil {
		fmt.Printf("WARNING: event bus not configured, [%s] for <%s/%s> not published\n", eventName, submissionId, bagId)
		return nil
	}
	detail, _ := workflowPayload(submissionId, bagId, extra)
	ev := uvaaptsbus.UvaBusEvent{
		EventName: eventName,
		ClientId:  clientId,
		Detail:    detail,
	}

	fmt.Printf("INFO: publishing [%v]\n", ev)
	err := bus.PublishEvent(&ev)
	if err != nil {
		fmt.Printf("ERROR: publishing (%s)\n", err.Error())
	}
	return err
}

func workflowPayload(submissionId string, bagId string, extra string) (json.RawMessage, error) {
	pl := uvaaptsbus.UvaWorkflowEvent{SubmissionId: submissionId, BagId: bagId, Extra: extra}
	return pl.Serialize()
}

//
// end of file
//
