package main

type AptStatusResponse struct {
	Count  int      `json:"count"`  // the client identifier
	Status []string `json:"status"` // the bags to be included in this submission
}

//
// end of file
//
