package main

import (
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// submission needs to be validated
func handleSubmissionValidate(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusValidating)
}

func handleSubmissionValidateFail(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusError)
}

func handleSubmissionReconcileFail(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusError)
}

// submission needs to be approved
func handleSubmissionApprove(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusPendingApproval)
}

// submission was approved
func handleSubmissionApproved(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// audit the approval cos the approver is contained in this event
	err := dao.AddApproval(workflowEvent.SubmissionId, workflowEvent.Extra)
	if err != nil {
		return err
	}

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusBuilding)
}

//
// end of file
//
