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
	"github.com/rs/xid"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

type Request struct {
	ClientIdentifier string `json:"cid"`        // the client identifier
	Collection       string `json:"collection"` // the collection name for the submission (optional)
	Storage          string `json:"storage"`    // the APT storage to use for this submission (optional)
}

type Response struct {
	SubmissionIdentifier string `json:"sid"`
	DepositBucket        string `json:"bucket"`
	DepositPath          string `json:"path"`
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	fmt.Printf("DEBUG: request [%s]\n", request.Body)
	r := Request{}
	err := json.Unmarshal([]byte(request.Body), &r)
	if err != nil {
		fmt.Printf("ERROR: json.Unmarshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// ensure we have the parameters we need
	if len(r.ClientIdentifier) == 0 {
		fmt.Printf("ERROR: missing required params (cid)\n")
		err = fmt.Errorf("missing required params (cid)")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

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
	c, err := dao.GetClientByIdentifier(r.ClientIdentifier)
	if err != nil {
		if errors.As(err, &uvaaptsdao.ErrClientNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusForbidden, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// create new submission identifier
	sid := newSubmissionIdentifier()

	// deal with default behavior
	if len(r.Collection) == 0 {
		r.Collection = sid
	}

	if len(r.Storage) == 0 {
		r.Storage = c.DefaultStorage
	}

	// create the new submission
	err = dao.AddSubmission(sid, c.Identifier, r.Collection, r.Storage)
	if err != nil {
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	fmt.Printf("DEBUG: new submission [%s], collection name [%s], storage [%s]\n", sid, r.Collection, r.Storage)

	// construct the response
	response := Response{}
	response.SubmissionIdentifier = sid
	response.DepositBucket = cfg.InboundBucket
	// S3 assets in <bucket>/<clientId>/<submissionId>/...
	response.DepositPath = fmt.Sprintf("%s/%s", r.ClientIdentifier, sid)

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	fmt.Printf("DEBUG: response [%s]\n", string(buf))
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

func newSubmissionIdentifier() string {
	return fmt.Sprintf("sid-%s", xid.New().String())
}

//
// end of file
//
