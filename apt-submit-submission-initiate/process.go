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
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
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
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	fmt.Printf("DEBUG: request [%s]\n", request.Body)
	r := Request{}
	err := json.Unmarshal([]byte(request.Body), &r)
	if err != nil {
		fmt.Printf("ERROR: json.Unmarshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// just to make sure...
	if len(r.BagFolders) == 0 {
		err := fmt.Errorf("no bags specified in request body")
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
	cli, err := dao.GetClientByIdentifier(cid)
	if err != nil {
		if errors.As(err, &ErrClientNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusForbidden, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// get the submission
	sub, err := dao.GetSubmissionByIdentifier(sid)
	if err != nil {
		if errors.As(err, &ErrSubmissionNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusNotFound, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// create our s3 helper client
	s3Client, err := newS3Client()
	if err != nil {
		fmt.Printf("ERROR: creating s3 client (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// S3 assets in <bucket>/<clientId>/<submissionId>/...
	submissionKeyPrefix := fmt.Sprintf("%s/%s", cli.Name, sub.Identifier)

	// get a complete list of all the files included in the specified submission
	suppliedFiles, err := s3Client.s3List(cfg.InboundBucket, submissionKeyPrefix)
	if err != nil {
		fmt.Printf("ERROR: listing submission assets (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// get all the bags included in the submission
	bagList := findIncludedBags(submissionKeyPrefix, suppliedFiles)
	if len(bagList) == 0 {
		err = fmt.Errorf("no bags included in the submission")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	// ensure the bags specified in the request are the same as the ones located
	if areIdentical(bagList, r.BagFolders) == false {
		err = fmt.Errorf("bag list does not match submission")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	// get an enumeration of all the files specified in the manifests
	itemizedFiles := make([]ManifestRow, 0)
	for _, bag := range bagList {
		rows, err := manifestContents(s3Client, cfg.InboundBucket, submissionKeyPrefix, bag)
		if err != nil {
			err = fmt.Errorf("manifest %s/%s bad or missing", bag, manifestName)
			return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
		}
		itemizedFiles = append(itemizedFiles, rows...)
	}

	// our enumerated files and the supplied list should be the same size
	if len(itemizedFiles)+len(bagList) != len(suppliedFiles) {
		err = fmt.Errorf("manifests do not match submission")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	// FIXME
	// ensure the enumerated list matches the provided list
	//

	// create our event bus client
	eventBus, _ := NewEventBus(cfg.BusName, cfg.BusEventSource)

	// create the bags
	err = createDBBags(dao, bagList, sid)
	if err != nil {
		fmt.Printf("ERROR: creating bags (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// create the files
	err = createDBFiles(dao, itemizedFiles, sid)
	if err != nil {
		fmt.Printf("ERROR: creating files (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// construct the response
	response := Response{}

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// we are done, publish the appropriate event and terminate
	_ = publishWorkflowEvent(eventBus, uvaaptsbus.EventSubmissionValidate, cli.Name, sid, "")

	fmt.Printf("DEBUG: response [%s]\n", string(buf))
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

//
// end of file
//
