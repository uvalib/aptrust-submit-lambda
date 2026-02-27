//
//
//

package main

import "time"

type Client struct {
	Name       string    `json:"name"`       // client name
	Identifier string    `json:"identifier"` // client identifier
	Created    time.Time `json:"created"`    // created time
}

type Submission struct {
	Identifier string    `json:"identifier"` // submission identifier
	Client     string    `json:"client"`     // owning client
	Created    time.Time `json:"created"`    // created time
}

type Bag struct {
	Name       string    `json:"name"`       // bag name
	Identifier string    `json:"identifier"` // bag identifier
	Submission string    `json:"submission"` // owning submission
	Created    time.Time `json:"created"`    // created time
}

type SubmissionStatus struct {
	Identifier string    `json:"identifier"` // submission identifier
	Status     string    `json:"status"`     // current status
	Updated    time.Time `json:"updated"`    // created time
}

type BagStatus struct {
	Identifier string    `json:"identifier"` // submission identifier
	Status     string    `json:"status"`     // current status
	Updated    time.Time `json:"updated"`    // created time
}

//
// end of file
//
