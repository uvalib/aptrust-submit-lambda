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
)

type Request struct {
	BagFolders []string `json:"bag_folders"`
}

type Response struct {
	//Sid string `json:"sid"`
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

	fmt.Printf("DEBUG: request [%s]\n", request.Body)
	r := Request{}
	err := json.Unmarshal([]byte(request.Body), &r)
	if err != nil {
		fmt.Printf("ERROR: json.Unmarshal() failed (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
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

	// get the client details
	_, err = dao.GetClientByIdentifier(cid)
	if err != nil {
		if errors.Is(err, ErrClientNotFound) {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusForbidden}, err
		}
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// get the submission
	_, err = dao.GetSubmissionByIdentifier(sid)
	if err != nil {
		if errors.Is(err, ErrSubmissionNotFound) {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusNotFound}, err
		}
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// do more stuff

	// construct the response
	response := Response{}
	//response.Sid = "sid-xx-example-xx"

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
