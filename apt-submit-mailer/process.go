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
	case uvaaptsbus.EventSubmissionApprove:
	case uvaaptsbus.EventSubmissionValidateFail:
	case uvaaptsbus.EventSubmissionReconcileFail:
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

	// assume the recipient is the administrator
	recipient := cfg.AdminEmail

	// if this is (potentially) an approval email, ensure this is not an auto-approve submission
	if be.EventName == uvaaptsbus.EventSubmissionApprove {
		// get the submission
		sub, err := dao.GetSubmissionByIdentifier(wf.SubmissionId)
		if err != nil {
			return err
		}

		// this will be an auto-approval so abandon the mailer
		if len(sub.ApprovalEmail) == 0 {
			fmt.Printf("INFO: auto-approve submission, no email necessary\n")
			return nil
		}

		// recipient will be the submission approver
		recipient = sub.ApprovalEmail
	}

	// render the email body
	subject, body, err := renderSubjectAndBody(cfg, recipient, be, wf)
	if err != nil {
		fmt.Printf("ERROR: rendering email content (%s)\n", err.Error())
		return err
	}

	// send the mail
	return sendEmail(cfg, subject, recipient, cfg.MailCC, body)
}

//
// end of file
//
