//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiResponse struct {
	Sid string `json:"sid"`
	// other stuff
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var cid string
	var sid string

	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("DEBUG: query param [%s] = [%s]\n", key, value)
		switch key {
		case "cid":
			cid = value
		case "sid":
			sid = value
		}
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	// ensure we have the parameters we need
	if len(cid) == 0 || len(sid) == 0 {
		err := fmt.Errorf("missing required query params: [cid, sid]")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, err
	}

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// create the data access object
	dao, err := newDao(cfg)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// cleanup on exit
	defer dao.Close()

	// construct the response
	response := ApiResponse{}
	response.Sid = "sid-xx-example-xx"

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}
	fmt.Printf("DEBUG: response [%s]\n", string(buf))
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

//
// end of file
//
