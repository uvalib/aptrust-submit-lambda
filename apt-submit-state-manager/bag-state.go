package main

import (
	"encoding/json"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// bag was built
func handleBagBuilt(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	return dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, uvaaptsdao.BagStatusReady)
}

// bag was submitted to APT
func handleBagSubmitted(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// apply the etag cos it is contained in this event
	extra := BagSubmittedEventExtraPayload{}
	_ = json.Unmarshal([]byte(workflowEvent.Extra), &extra)

	err := dao.UpdateBagETag(workflowEvent.BagId, workflowEvent.SubmissionId, extra.ETag)
	if err != nil {
		return err
	}

	// update the state of the bag
	err = dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, uvaaptsdao.BagStatusPendingIngest)
	if err != nil {
		return err
	}

	// get the submission status
	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}

	// if the status is 'building' update to 'pending-ingest'
	if ss.State == uvaaptsdao.SubmissionStatusBuilding {
		// update the status of the submission
		err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusPendingIngest)
	}
	return err
}

// bag was rejected by APT
func handleBagRejected(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// update the state of the bag
	err := dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, uvaaptsdao.BagStatusError)
	if err != nil {
		return err
	}

	// get the submission status
	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}

	// if the status is 'pending-ingest', update to 'incomplete'
	if ss.State == uvaaptsdao.SubmissionStatusPendingIngest {
		// update the status of the submission
		err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusIncomplete)
	}
	return err
}

// bag was successfully accepted by APT
func handleBagAccepted(bus uvaaptsbus.UvaBus, busEvent *uvaaptsbus.UvaBusEvent, workflowEvent *uvaaptsbus.UvaWorkflowEvent, dao *uvaaptsdao.Dao) error {

	// in addition to updating the just accepted bag state
	// check to see if a) the current submission status is in the 'pending-ingest' state and
	// b) if the submission bag(s) are done
	//

	// get the submission status
	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}

	// if the status is 'pending-ingest', check to see if all the bags are done
	if ss.State == uvaaptsdao.SubmissionStatusPendingIngest {

		// check to see how many bags remain in the pending state... if none
		// update the submission state to 'complete'
		bags, err := dao.GetBagsBySubmission(workflowEvent.SubmissionId)
		if err != nil {
			return err
		}
		allDone := true
		for _, b := range bags {
			if b.Name != workflowEvent.BagId {
				bs, err := dao.GetBagStateBySubmissionAndName(workflowEvent.SubmissionId, b.Name)
				if err != nil {
					return err
				}
				if bs.State == uvaaptsdao.SubmissionStatusPendingIngest {
					allDone = false
					break
				}
			}
		}
		if allDone == true {
			// update the status of the submission
			err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusComplete)
			if err != nil {
				return err
			}
			err = publishWorkflowEvent(bus, uvaaptsbus.EventSubmissionComplete, busEvent.ClientId, workflowEvent.SubmissionId, "", "")
			if err != nil {
				return err
			}
		}
	}

	// update the state of the bag
	return dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, uvaaptsdao.BagStatusComplete)
}

//
// end of file
//
