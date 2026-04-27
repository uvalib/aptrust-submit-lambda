//
// main message processing
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"

	"math/rand"
)

// after this period we will log a warning about a stalled ingest
var stalledPendingThreshold = 24 * time.Hour

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
	bags, err := dao.GetBagsByStatus(uvaaptsdao.BagStatusPendingIngest)
	if err != nil {
		if errors.As(err, &uvaaptsdao.ErrBagNotFound) {
			fmt.Printf("INFO: no bags in '%s' status\n", uvaaptsdao.BagStatusPendingIngest)
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

		// if we have more bags than requests we are permitted...
		if len(bags) > cfg.MaxRequests {
			fmt.Printf("INFO: randomly selecting the first %d of %d bags\n", cfg.MaxRequests, len(bags))
			rand.Shuffle(len(bags), func(i, j int) {
				bags[i], bags[j] = bags[j], bags[i]
			})

			// Take the first n elements
			bags = bags[:cfg.MaxRequests]
		} else {
			fmt.Printf("INFO: checking status on %d bags\n", len(bags))
		}

		// proces each of the bags we know about
		for _, bg := range bags {

			if len(bg.ETag) != 0 {
				// get the status, ignore errors
				status, _ := getAptStatus(cfg, httpClient, bg)
				switch status {
				// interesting statuses that generate events
				case AptStatusCancelled, AptStatusFailed, AptStatusSuspended, AptStatusSuccess:

					// get the submission details for the bag
					sub, err := dao.GetSubmissionByIdentifier(bg.Submission)
					if err != nil {
						return err
					}

					if status == AptStatusSuccess {
						// victory, fire the accepted event
						_ = publishWorkflowEvent(eventBus, uvaaptsbus.EventBagAccepted, sub.Client, bg.Submission, bg.Name, "")
					} else {
						// something terminal happened, fire the rejected event
						_ = publishWorkflowEvent(eventBus, uvaaptsbus.EventBagRejected, sub.Client, bg.Submission, bg.Name, "")
					}

				case AptStatusPending:

					// get the latest bag status to determine when we submitted to APT and
					// log a message if it is stale (stuck, often happens)
					bstat, err := dao.GetBagStateBySubmissionAndName(bg.Submission, bg.Name)
					if err != nil {
						return err
					}
					age := time.Since(bstat.Updated)
					if age > stalledPendingThreshold {
						fmt.Printf("WARNING: ingest stalled for <%s/%s> (submitted: %s)\n", bg.Submission, bg.Name, bstat.Updated)
					}

				case AptStatusStarted:
					// do nothing

				case AptStatusUnknown:
					// do nothing

				default:
					fmt.Printf("ERROR: unexpected status for bag <%s/%s> (%s)\n", bg.Submission, bg.Name, status)
				}
			} else {
				fmt.Printf("ERROR: bag <%s/%s> has an empty etag, cannot check for status\n", bg.Submission, bg.Name)
			}
		}
	}

	return nil
}

//
// end of file
//
