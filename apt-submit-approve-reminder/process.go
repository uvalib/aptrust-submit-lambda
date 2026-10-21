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

	// get the submissions
	submissions, err := dao.GetSubmissionsByStatus(uvaaptsdao.SubmissionStatusPendingApproval)
	if err != nil {
		if errors.Is(err, uvaaptsdao.ErrSubmissionNotFound) {
			fmt.Printf("INFO: no submissions in '%s' status\n", uvaaptsdao.SubmissionStatusPendingApproval)
			return nil
		}
		return err
	}

	if len(submissions) != 0 {
		fmt.Printf("INFO: %d submission(s) in '%s' status\n", len(submissions), uvaaptsdao.SubmissionStatusPendingApproval)

		// create our event bus client
		eventBus, _ := NewEventBus(cfg.BusName, cfg.BusEventSource)

		// proces each of the submissions we know about
		for _, s := range submissions {

			// get the number of days since the submission went into the pending approval state
			duration := time.Since(s.Created)
			daysSince := int(duration.Hours() / 24)

			createdYYYYMMDD := s.Created.Format("2006-01-02")

			// if approval was originally issued today
			if daysSince == 0 {
				fmt.Printf("INFO: submission %s, original approval mail sent %s, no reminder required\n", s.Identifier, createdYYYYMMDD)
			} else {
				// how many days before we send a reminder
				remaining := daysSince % cfg.ReminderDays
				if remaining == 0 {
					fmt.Printf("INFO: submission %s, original approval mail sent %s, sending reminder email\n", s.Identifier, createdYYYYMMDD)
					err = publishWorkflowEvent(eventBus,
						uvaaptsbus.EventCommandSubmissionEmailApprove,
						s.Client,
						s.Identifier,
						"", "")
				} else {
					fmt.Printf("INFO: submission %s, original approval mail sent %s, reminder email in %d day(s)\n", s.Identifier, createdYYYYMMDD, remaining)
				}
			}
		}
	}

	return nil
}

//
// end of file
//
