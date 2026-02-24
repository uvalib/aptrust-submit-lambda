//
//
//

package main

import "time"

type Submission struct {
	Id         int       `json:"id"`         // internal submission id, opaque
	Identifier string    `json:"identifier"` // submission identifier
	ClientId   int       `json:"client_id"`  // owning client
	Created    time.Time `json:"created"`    // created time
}

type Client struct {
	Id         int       `json:"id"`         // internal client id, opaque
	Name       string    `json:"name"`       // client name
	Identifier string    `json:"identifier"` // client identifier
	Created    time.Time `json:"created"`    // created time
}

type SubmissionStatus struct {
	Identifier string `json:"identifier"` // submission identifier
	Status     string `json:"status"`     // current status
}

//
// end of file
//
