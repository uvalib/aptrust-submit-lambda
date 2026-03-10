//
// main message processing
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

type Response struct {
	Submission string    `json:"submission"`
	Status     string    `json:"status"`
	Updated    time.Time `json:"updated"`
	// other stuff
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var sid string

	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("DEBUG: query param [%s] = [%s]\n", key, value)
		switch key {
		case "sid":
			sid = value
		}
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	// ensure we have the parameters we need
	if len(sid) == 0 {
		err := fmt.Errorf("missing required query params: [sid]")
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

	// get the submission
	s, err := dao.GetSubmissionByIdentifier(sid)
	if err != nil {
		if errors.As(err, &ErrSubmissionNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusNotFound, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// construct the response
	response := Response{}
	response.Submission = s.Identifier
	response.Status = s.Status
	response.Updated = s.Updated

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
