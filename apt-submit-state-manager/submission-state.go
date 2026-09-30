package main

import (
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// an event asking for a transition the submission cannot make is not worth retrying: the
// state was just read from the database and a redelivery will read the same value. Report
// whether the caller should proceed rather than returning an error, so a duplicate or out
// of order event is discarded instead of exhausting its retries and landing in the DLQ
func submissionCanTransition(submissionId string, current string, required string, target string) bool {

	switch current {
	case required:
		return true

	case target:
		// the event has already been applied, most likely a redelivery
		fmt.Printf("INFO: submission [%s] is already '%s', ignoring\n", submissionId, target)

	default:
		fmt.Printf("WARNING: submission [%s] in incorrect state for '%s' (%s), ignoring\n", submissionId, target, current)
	}

	return false
}

func handleSubmissionReconcileFail(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of all the bags
	bags, err := dao.GetBagsBySubmission(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}
	for _, b := range bags {
		err = dao.UpdateBagState(b.Name, workflowEvent.SubmissionId, uvaaptsdao.BagStatusError)
		if err != nil {
			return err
		}
	}

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusError)
}

// submission was abandoned
func handleSubmissionAbandoned(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		fmt.Printf("ERROR: getting submission state (%s)\n", err.Error())
		return err
	}

	// validate that the submission state is as expected
	if submissionCanTransition(workflowEvent.SubmissionId, ss.State, uvaaptsdao.SubmissionStatusPendingApproval, uvaaptsdao.SubmissionStatusAbandoned) == false {
		return nil
	}

	// update the state of all the bags
	bags, err := dao.GetBagsBySubmission(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}
	for _, b := range bags {
		err = dao.UpdateBagState(b.Name, workflowEvent.SubmissionId, uvaaptsdao.BagStatusAbandoned)
		if err != nil {
			return err
		}
	}

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusAbandoned)
}

func handleSubmissionIncomplete(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		fmt.Printf("ERROR: getting submission state (%s)\n", err.Error())
		return err
	}

	// validate that the submission state is as expected
	if submissionCanTransition(workflowEvent.SubmissionId, ss.State, uvaaptsdao.SubmissionStatusError, uvaaptsdao.SubmissionStatusIncomplete) == false {
		return nil
	}

	// update the state of all the bags
	//bags, err := dao.GetBagsBySubmission(workflowEvent.SubmissionId)
	//if err != nil {
	//	return err
	//}
	//for _, b := range bags {
	//	err = dao.UpdateBagState(b.Name, workflowEvent.SubmissionId, uvaaptsdao.BagStatusAbandoned)
	//	if err != nil {
	//		return err
	//	}
	//}

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusIncomplete)
}

//
// end of file
//
