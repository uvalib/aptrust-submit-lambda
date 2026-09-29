package main

import (
	"encoding/json"
	"fmt"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// a bag in one of these states has finished its journey and will not transition
// again without operator intervention
func bagStateIsTerminal(state string) bool {
	switch state {
	case uvaaptsdao.BagStatusComplete,
		uvaaptsdao.BagStatusError,
		uvaaptsdao.BagStatusAbandoned:
		return true
	}
	return false
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

	// update the state of the accepted bag first; the completeness check below then
	// sees a consistent view of every bag and does not need to special case this one
	err := dao.UpdateBagState(workflowEvent.BagId, workflowEvent.SubmissionId, uvaaptsdao.BagStatusComplete)
	if err != nil {
		return err
	}

	// get the submission status
	ss, err := dao.GetSubmissionStateByIdentifier(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}

	// only a submission that is awaiting ingest can transition to complete; if a bag was
	// rejected the submission is already 'incomplete' and must stay that way
	if ss.State != uvaaptsdao.SubmissionStatusPendingIngest {
		return nil
	}

	// the submission is complete only when every one of its bags has reached a
	// terminal state; a bag that is still registered/building/ready/submitting or
	// awaiting ingest means there is more work to come
	bags, err := dao.GetBagsBySubmission(workflowEvent.SubmissionId)
	if err != nil {
		return err
	}
	for _, b := range bags {
		bs, err := dao.GetBagStateBySubmissionAndName(workflowEvent.SubmissionId, b.Name)
		if err != nil {
			return err
		}
		if bagStateIsTerminal(bs.State) == false {
			fmt.Printf("INFO: submission <%s> not complete, bag <%s> is '%s'\n", workflowEvent.SubmissionId, b.Name, bs.State)
			return nil
		}
	}

	// update the status of the submission
	err = dao.UpdateSubmissionState(workflowEvent.SubmissionId, uvaaptsdao.SubmissionStatusComplete)
	if err != nil {
		return err
	}

	return publishWorkflowEvent(bus, uvaaptsbus.EventSubmissionComplete, busEvent.ClientId, workflowEvent.SubmissionId, "", "")
}

//
// end of file
//
