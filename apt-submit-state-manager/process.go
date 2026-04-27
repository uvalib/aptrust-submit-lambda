//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to uvaaptsbus event
	be, err := uvaaptsbus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling bus event (%s)\n", err.Error())
		return err
	}

	// ensure this is the type of event we want to process
	switch be.EventName {
	case uvaaptsbus.EventSubmissionValidate: // submission ready to validate
	case uvaaptsbus.EventSubmissionValidateFail: // validation failed
	//case uvaaptsbus.EventSubmissionReconcile: // submission ready to reconcile
	case uvaaptsbus.EventSubmissionReconcileFail: // reconciliation failed
	case uvaaptsbus.EventSubmissionApprove: // submission to be approved
	case uvaaptsbus.EventSubmissionApproved: // submission approved (or auto approved)
	case uvaaptsbus.EventSubmissionAbandoned: // submission abandoned

	case uvaaptsbus.EventBagBuilt: // bag built
	case uvaaptsbus.EventBagSubmitted: // bag submitted
	case uvaaptsbus.EventBagRejected: // bag rejected (by APT)
	case uvaaptsbus.EventBagAccepted: // bag accepted (by APT)
	default:
		fmt.Printf("WARNING: unexpected event type (%s), ignoring\n", be.EventName)
		return nil
	}

	// make the workflow event
	wf, err := uvaaptsbus.MakeWorkflowEvent(be.Detail)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling workflow event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: event %s/%s\n", be.String(), wf.String())

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// create the data access object
	dao, err := uvaaptsdao.NewDao(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	if err != nil {
		return err
	}

	// cleanup on exit
	defer dao.Close()

	// create our event bus client
	eventBus, _ := NewEventBus(cfg.BusName, cfg.BusEventSource)

	// event specific processing
	switch be.EventName {
	//
	// submission events
	//
	case uvaaptsbus.EventSubmissionValidate:
		err = handleSubmissionValidate(eventBus, be, wf, dao)

	case uvaaptsbus.EventSubmissionValidateFail:
		err = handleSubmissionValidateFail(eventBus, be, wf, dao)

	//case uvaaptsbus.EventSubmissionReconcile: err = handleSubmissionReconcile(eventBus, be, wf, dao)
	case uvaaptsbus.EventSubmissionReconcileFail:
		err = handleSubmissionReconcileFail(eventBus, be, wf, dao)

	case uvaaptsbus.EventSubmissionApprove:
		err = handleSubmissionApprove(eventBus, be, wf, dao)

	case uvaaptsbus.EventSubmissionApproved:
		err = handleSubmissionApproved(eventBus, be, wf, dao)

	case uvaaptsbus.EventSubmissionAbandoned:
		err = handleSubmissionAbandoned(eventBus, be, wf, dao)

	//
	// bag events
	//
	case uvaaptsbus.EventBagBuilt:
		err = handleBagBuilt(eventBus, be, wf, dao)

	case uvaaptsbus.EventBagSubmitted:
		err = handleBagSubmitted(eventBus, be, wf, dao)

	case uvaaptsbus.EventBagRejected:
		err = handleBagRejected(eventBus, be, wf, dao)

	case uvaaptsbus.EventBagAccepted:
		err = handleBagAccepted(eventBus, be, wf, dao)
	}

	return err
}

//
// end of file
//
