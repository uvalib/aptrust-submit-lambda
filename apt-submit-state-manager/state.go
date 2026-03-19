package main

import (
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// submission needs to be approved
func handleSubmissionApprove(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusPendingApproval)
}

// submission was approved
func handleSubmissionApproved(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the submission
	return dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusBuilding)
}

// bag was submitted to APT
func handleBagSubmitted(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	err := dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, BagStatusPendingIngest)
	if err != nil {
		return err
	}

	// also apply the etag cos it is contained in this event
	err = dao.UpdateBagETag(workflowEvent.BagId, workflowEvent.SubmissionId, workflowEvent.Extra)

	return err
}

// bag was rejected by APT
func handleBagRejected(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	return dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, BagStatusError)
}

// bag was successfully accepted by APT
func handleBagAccepted(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	return dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, bagStatusComplete)
}

//
// end of file
//
