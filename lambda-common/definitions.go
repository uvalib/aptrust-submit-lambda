//
//
//

package main

import "fmt"

// error definitions
var ErrClientNotFound = fmt.Errorf("client not found")
var ErrSubmissionNotFound = fmt.Errorf("submission not found")

// submission status definitions
var SubmissionStatusRegistered = "registered"
var SubmissionStatusValidating = "validating"
var SubmissionStatusBuilding = "building"
var SubmissionStatusPendingApproval = "pending-approval"
var SubmissionStatusSubmitting = "submitting"
var SubmissionStatusPendingIngest = "pending-ingest"
var SubmissionStatusError = "error"
var SubmissionStatusComplete = "complete"

//
// end of file
//
