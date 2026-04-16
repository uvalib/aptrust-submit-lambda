//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to uvaaptsbus event
	be, err := uvaaptsbus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling bus event (%s)\n", err.Error())
		return err
	}

	// ensure this is the type of event we want to process
	switch be.EventName {
	case uvaaptsbus.EventSubmissionComplete: // submission complete and ready to purge
	case uvaaptsbus.EventCommandSubmissionPurge: // submission complete and ready to purge
	case uvaaptsbus.EventCommandBagPurge: // bag complete and ready to purge
	default:
		fmt.Printf("WARNING: unexpected event type (%s), ignoring\n", be.EventName)
		return nil
	}

	// make the workflow event
	wf, err := uvaaptsbus.MakeWorkflowEvent(be.Detail)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling workflow event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: event %s/%s\n", be.String(), wf.String())

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// create our s3 helper client
	s3Client, err := newS3Client()
	if err != nil {
		fmt.Printf("ERROR: creating s3 client (%s)\n", err.Error())
		return err
	}

	// event specific processing
	switch be.EventName {
	//
	// submission events
	//
	case uvaaptsbus.EventSubmissionComplete, uvaaptsbus.EventCommandSubmissionPurge:
		//err = handleSubmissionValidate(eventBus, be, wf, dao)

	case uvaaptsbus.EventCommandBagPurge:
		//err = handleSubmissionValidateFail(eventBus, be, wf, dao)
	}

	// S3 assets in <bucket>/<clientId>/<submissionId>/...
	pathPrefix := path.Join(be.ClientId, wf.SubmissionId)

	// get a complete list of all the files included in the specified submission
	s3files, err := s3Client.s3List(cfg.AssetBucket, pathPrefix)
	if err != nil {
		fmt.Printf("ERROR: listing S3 assets (%s)\n", err.Error())
		//return err
	}

	fmt.Printf("INFO: located %d assets in [s3://%s/%s]\n", len(s3files), cfg.AssetBucket, pathPrefix)

	dir := path.Join(cfg.AssetFilesystem, pathPrefix)
	de, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("ERROR: reading [%s] (%s)\n", dir, err.Error())
		//return err
	}

	fmt.Printf("INFO: located %d assets in [%s]\n", len(de), dir)

	for _, d := range de {
		//err = os.RemoveAll(path.Join(dir, d.Name()))
		fmt.Printf("INFO: %s\n", path.Join(dir, d.Name()))
		//if err != nil {
		//	return err
		//}
	}

	return err
}

//
// end of file
//
