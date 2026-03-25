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
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

type Request struct {
	ClientIdentifier     string   `json:"cid"`         // the client identifier
	SubmissionIdentifier string   `json:"sid"`         // the submission identifier
	BagFolders           []string `json:"bag_folders"` // the bags to be included in this submission
}

type Response struct {
	Submission string    `json:"submission"`
	Status     string    `json:"status"`
	Updated    time.Time `json:"updated"`
	// other stuff
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
	if len(r.ClientIdentifier) == 0 || len(r.SubmissionIdentifier) == 0 || len(r.BagFolders) == 0 {
		fmt.Printf("ERROR: one or more missing required params: [cid, sid, bag_folders]\n")
		err = fmt.Errorf("one or more missing required params: [cid, sid, bag_folders]")
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
	cli, err := dao.GetClientByIdentifier(r.ClientIdentifier)
	if err != nil {
		if errors.As(err, &ErrClientNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusForbidden, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// get the submission
	sub, err := dao.GetSubmissionByIdentifier(r.SubmissionIdentifier)
	if err != nil {
		if errors.As(err, &ErrSubmissionNotFound) {
			return apiGatewayProxyErrorResponse(http.StatusNotFound, err)
		}
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// ensure this submission belongs to this client
	if sub.Client != cli.Identifier {
		fmt.Printf("ERROR: client does not match submission identifier\n")
		err = fmt.Errorf("client does not match submission identifier")
		return apiGatewayProxyErrorResponse(http.StatusForbidden, err)
	}

	fmt.Printf("DEBUG: processing submission of %d bags\n", len(r.BagFolders))

	// create our s3 helper client
	s3Client, err := newS3Client()
	if err != nil {
		fmt.Printf("ERROR: creating s3 client (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// S3 assets in <bucket>/<clientId>/<submissionId>/...
	submissionKeyPrefix := fmt.Sprintf("%s/%s", cli.Identifier, sub.Identifier)

	// get a complete list of all the files included in the specified submission
	suppliedFiles, err := s3Client.s3List(cfg.InboundBucket, submissionKeyPrefix)
	if err != nil {
		fmt.Printf("ERROR: listing submission assets (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// get all the bags included in the submission
	bagList := findIncludedBags(submissionKeyPrefix, suppliedFiles)
	if len(bagList) == 0 {
		fmt.Printf("ERROR: no bags included in the submission\n")
		err = fmt.Errorf("no bags included in the submission")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	// ensure the bags specified in the request are the same as the ones located
	if areIdentical(bagList, r.BagFolders) == false {
		fmt.Printf("ERROR: located bag list does not match submission list\n")
		err = fmt.Errorf("located bag list does not match submission list")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	// get an enumeration of all the files specified in the manifests
	itemizedFiles := make([]ManifestRow, 0)
	for _, bag := range bagList {
		rows, err := manifestContents(s3Client, cfg.InboundBucket, submissionKeyPrefix, bag)
		if err != nil {
			fmt.Printf("ERROR: manifest %s/%s bad or missing\n", bag, manifestName)
			err = fmt.Errorf("manifest %s/%s bad or missing", bag, manifestName)
			return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
		}
		itemizedFiles = append(itemizedFiles, rows...)
	}

	fmt.Printf("INFO: %d files enumerated in %d manifests\n", len(itemizedFiles), len(bagList))
	fmt.Printf("INFO: %d files located in the submission\n", len(suppliedFiles))

	// our enumerated files and the supplied list should be the same size
	if len(itemizedFiles)+len(bagList) != len(suppliedFiles) {
		fmt.Printf("ERROR: manifests do not match submission\n")
		err = fmt.Errorf("manifests do not match submission")
		return apiGatewayProxyErrorResponse(http.StatusBadRequest, err)
	}

	// FIXME
	// ensure the enumerated list matches the provided list
	//

	// create our event bus client
	eventBus, _ := NewEventBus(cfg.BusName, cfg.BusEventSource)

	// create the bags
	err = createDBBags(dao, bagList, r.SubmissionIdentifier)
	if err != nil {
		fmt.Printf("ERROR: creating bags (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// create the files
	err = createDBFiles(dao, itemizedFiles, r.SubmissionIdentifier)
	if err != nil {
		fmt.Printf("ERROR: creating files (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// construct the response
	response := Response{}
	response.Submission = r.SubmissionIdentifier
	response.Status = SubmissionStatusValidating
	response.Updated = time.Now().UTC()

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return apiGatewayProxyErrorResponse(http.StatusInternalServerError, err)
	}

	// we are done, publish the appropriate event and terminate
	_ = publishWorkflowEvent(eventBus, uvaaptsbus.EventSubmissionValidate, cli.Identifier, r.SubmissionIdentifier, "", "")

	fmt.Printf("DEBUG: response [%s]\n", string(buf))
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

//
// end of file
//
