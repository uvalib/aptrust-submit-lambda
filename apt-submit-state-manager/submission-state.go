package main

import (
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// submission needs to be validated
func handleSubmissionValidate(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusValidating)
}

func handleSubmissionValidateFail(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusError)
}

// submission needs to be reconciled
//func handleSubmissionReconcile(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {
//}

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
	if ss.State != uvaaptsdao.SubmissionStatusPendingApproval {
		err = fmt.Errorf("submission [%s] in incorrect state for abandon (%s)", workflowEvent.SubmissionId, ss.State)
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
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

//
// end of file
//
