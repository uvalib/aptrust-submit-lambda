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
	Name       string    `json:"name"`
	State      string    `json:"state"`
	Updated    time.Time `json:"updated"`
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var bid string

	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("DEBUG: query param [%s] = [%s]\n", key, value)
		switch key {
		case "bid":
			bid = value
		}
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	// ensure we have the parameters we need
	if len(bid) == 0 {
		err := fmt.Errorf("missing required query params: [bid]")
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

	// get the bag state
	bs, err := dao.GetBagStateByName(bid)
	if err != nil {
		if errors.As(err, &uvaaptsdao.ErrBagNotFound) {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusNotFound}, nil
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// construct the response
	response := Response{}
	response.Submission = bs.Submission
	response.Name = bs.Name
	response.State = bs.State
	response.Updated = bs.Updated

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
