//
// main message processing
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

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

	// get the bags
	bags, err := dao.GetBagsByStatus(BagStatusPendingIngest)
	if err != nil {
		if errors.As(err, &ErrBagNotFound) {
			fmt.Printf("INFO: no bags in '%s' status\n", BagStatusPendingIngest)
			return nil
		}
		return err
	}

	if len(bags) != 0 {
		// create our event bus client
		eventBus, _ := NewEventBus(cfg.BusName, cfg.BusEventSource)

		// create our HTTP client
		httpClient := newHttpClient(1, cfg.HttpTimeout)
		// important, cleanup properly
		defer httpClient.CloseIdleConnections()

		// proces each of the bags we know about
		for _, bg := range bags {

			if len(bg.ETag) != 0 {
				// get the status, ignore errors
				status, _ := getAptStatus(cfg, httpClient, bg)
				switch status {
				case AptStatusCancelled:
				case AptStatusFailed:
				case AptStatusSuspended:
					// something terminal happened, fire the rejected event
					_ = publishWorkflowEvent(eventBus, uvaaptsbus.EventBagRejected, "", bg.Submission, bg.Name, "")

				case AptStatusSuccess:
					// victory, fire the accepted event
					_ = publishWorkflowEvent(eventBus, uvaaptsbus.EventBagAccepted, "", bg.Submission, bg.Name, "")

				default: // basically, do nothing
				}
			} else {
				fmt.Printf("WARNING: bag <%s/%s> has an empty etag, cannot check for status\n", bg.Submission, bg.Name)
			}
		}
	}

	return nil
}

//
// end of file
//
