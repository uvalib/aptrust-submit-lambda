package main

import (
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// bag was built
func handleBagBuilt(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	return dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, BagStatusReady)
}

// bag was submitted to APT
func handleBagSubmitted(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// apply the etag cos it is contained in this event
	err := dao.UpdateBagETag(workflowEvent.BagId, workflowEvent.SubmissionId, workflowEvent.Extra)
	if err != nil {
		return err
	}

	// update the state of the bag
	err = dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, BagStatusPendingIngest)
	if err != nil {
		return err
	}

	// get the submission status
	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}

	// if the status is 'building' update to 'pending-ingest'
	if ss.State == SubmissionStatusBuilding {
		// update the status of the submission
		err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusPendingIngest)
	}
	return err
}

// bag was rejected by APT
func handleBagRejected(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	err := dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, BagStatusError)
	if err != nil {
		return err
	}

	// get the submission status
	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}

	// if the status is 'pending-ingest', update to 'incomplete'
	if ss.State == SubmissionStatusPendingIngest {
		// update the status of the submission
		err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, SubmissionStatusIncomplete)
	}
	return err
}

// bag was successfully accepted by APT
func handleBagAccepted(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	return dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, bagStatusComplete)
}

//
// end of file
//
