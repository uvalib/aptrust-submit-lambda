//
//
//

package main

import "fmt"

// error definitions
var ErrClientNotFound = fmt.Errorf("client not found")
var ErrSubmissionNotFound = fmt.Errorf("submission not found")
var ErrBagNotFound = fmt.Errorf("bag not found")

// submission status definitions
var SubmissionStatusRegistered = "registered"
var SubmissionStatusValidating = "validating"
var SubmissionStatusBuilding = "building"
var SubmissionStatusPendingApproval = "pending-approval"
var SubmissionStatusSubmitting = "submitting"
var SubmissionStatusPendingIngest = "pending-ingest"
var SubmissionStatusError = "error"
var SubmissionStatusComplete = "complete"

// bag status definitions
var BagStatusBuilding = "building"
var BagStatusReady = "ready"
var BagStatusSubmitting = "submitting"
var BagStatusPendingIngest = "pending-ingest"
var BagStatusError = "error"
var bagStatusComplete = "complete"

//
// end of file
//
