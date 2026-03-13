//
// main message processing
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

type Request struct {
	Collection string `json:"collection"` // the collection name for the submission (optional)
	Storage    string `json:"storage"`    // the APT storage to use for this submission (optional)
}

type Response struct {
	SubmissionIdentifier string `json:"sid"`
	DepositBucket        string `json:"bucket"`
	DepositPath          string `json:"path"`
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var cid string

	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("DEBUG: query param [%s] = [%s]\n", key, value)
		switch key {
		case "cid":
			cid = value
		}
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	// ensure we have the parameters we need
	if len(cid) == 0 {
		err := fmt.Errorf("missing required query params: [cid]")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	fmt.Printf("DEBUG: request [%s]\n", request.Body)
	r := Request{}
	err := json.Unmarshal([]byte(request.Body), &r)
	if err != nil {
		fmt.Printf("ERROR: json.Unmarshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// FIXME
	fmt.Printf("DEBUG: collection [%s], storage [%s]\n", r.Collection, r.Storage)

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// create the data access object
	dao, err := uvaaptsdao.NewDao(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	if err != nil {
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// cleanup on exit
	defer dao.Close()

	// get the client details
	c, err := dao.GetClientByIdentifier(cid)
	if err != nil {
		if errors.As(err, &ErrClientNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusForbidden, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// create the new submission
	s, err := dao.CreateNewSubmission(c.Identifier)
	if err != nil {
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// construct the response
	response := Response{}
	response.SubmissionIdentifier = s.Identifier
	response.DepositBucket = cfg.InboundBucket
	// S3 assets in <bucket>/<clientId>/<submissionId>/...
	response.DepositPath = fmt.Sprintf("%s/%s", cid, s.Identifier)

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}
	fmt.Printf("DEBUG: response [%s]\n", string(buf))
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

//
// end of file
//
