package main

import (
	"encoding/json"
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// ready to approve the submission
func handleSubmissionApprove(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// get the client details
	cli, err := dao.GetClientByIdentifier(busEvent.ClientId)
	if err != nil {
		fmt.Printf("ERROR: client details not found [%s] (%s)\n", busEvent.ClientId, err.Error())
		return err
	}

	// if we do not require manual approval, just process as if it was approved
	// otherwise, it will be approved following user action
	if len(cli.ApprovalEmail) == 0 {
		// do the approval process
		err = submissionApproved(bus, dao, busEvent.ClientId, workflowEvent.SubmissionId, "<system>", "")
	} else {
		// update the state of the submission to reflect it is pending approval
		err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusPendingApproval)
		if err != nil {
			fmt.Printf("ERROR: updating submission state [%s] (%s)\n", workflowEvent.SubmissionId, err.Error())
			return err
		}
	}

	return err
}

// submission was approved (manually)
func handleSubmissionApproval(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		fmt.Printf("ERROR: getting submission state (%s)\n", err.Error())
		return err
	}

	// validate that the submission state is as expected
	if ss.State != uvaaptsdao.SubmissionStatusPendingApproval {
		err = fmt.Errorf("submission [%s] in incorrect state for approvals (%s)", workflowEvent.SubmissionId, ss.State)
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	// unpack the extra payload
	extra := ApprovalEventExtraPayload{}
	_ = json.Unmarshal([]byte(workflowEvent.Extra), &extra)

	// and do the approval process
	return submissionApproved(bus, dao, busEvent.ClientId, workflowEvent.SubmissionId, extra.ComputeID, extra.Storage)
}

func submissionApproved(bus uvaaptsbus.UvaBus, dao *uvaaptsdao.Dao, clientId string, submissionId string, approver string, storage string) error {

	// audit the approval
	err := dao.AddApproval(submissionId, approver)
	if err != nil {
		fmt.Printf("ERROR: adding approval record for [%s], continuing (%s)\n", submissionId, err.Error())
	}

	// update the storage for this submission
	if len(storage) != 0 {
		err = dao.UpdateSubmissionStorage(submissionId, storage)
		if err != nil {
			fmt.Printf("ERROR: updating storage for [%s], continuing (%s)\n", submissionId, err.Error())
			return err
		}
	}

	// get all the bags for this submission
	bags, err := dao.GetBagsBySubmission(submissionId)
	if err != nil {
		fmt.Printf("ERROR: getting submission bags (%s)\n", err.Error())
		return err
	}

	// and generate a bag initiate event for each one
	for _, bag := range bags {
		err = publishWorkflowEvent(bus, uvaaptsbus.EventBagInitiate, clientId, submissionId, bag.Name, "")
		if err != nil {
			fmt.Printf("ERROR: publishing bag initiate event for <%s:%s>, continuing (%s)\n", submissionId, bag.Name, err.Error())
		}
	}

	// update the state of the submission
	return dao.UpdateSubmissionState(submissionId, uvaaptsdao.SubmissionStatusBuilding)
}

//
// end of file
//
