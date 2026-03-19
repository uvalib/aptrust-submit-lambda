package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

type AptStatusResponse struct {
	Count   int         `json:"count"`
	Results []AptStatus `json:"results"`
}

type AptStatus struct {
	Status string `json:"status"`
}

// possible responses from the status endpoint
var AptStatusCancelled = "Cancelled"
var AptStatusFailed = "Failed"
var AptStatusPending = "Pending"
var AptStatusStarted = "Started"
var AptStatusSuspended = "Suspended"
var AptStatusSuccess = "Success"
var AptStatusUnknown = "Unknown"

func getAptStatus(cfg *Config, httpClient *http.Client, bag uvaaptsdao.Bag) (string, error) {

	fmt.Printf("DEBUG: checking for status of <%s/%s>\n", bag.Submission, bag.Name)

	// create the endpoint URL
	url := strings.Replace(cfg.APTStatusUrl, "{:etag}", bag.ETag, 1)
	b, err := httpGet(httpClient, url, cfg.APTUser, cfg.APTKey)
	if err == nil {
		res := AptStatusResponse{}
		err := json.Unmarshal(b, &res)
		if err == nil {
			if res.Count == 1 {
				fmt.Printf("INFO: bag status for <%s/%s> (%s)\n", bag.Submission, bag.Name, res.Results[0].Status)
				return res.Results[0].Status, nil
			} else {
				fmt.Printf("WARNING: cannot determine bag status for <%s/%s>\n", bag.Submission, bag.Name)
				return AptStatusUnknown, nil
			}
		} else {
			fmt.Printf("ERROR: json.Unmarshal() failed (%s)\n", err.Error())
		}
	} else {
		fmt.Printf("WARNING: getting bag status for <%s/%s> (%s)\n", bag.Submission, bag.Name, err.Error())
	}
	return AptStatusUnknown, err
}

//
// end of file
//
