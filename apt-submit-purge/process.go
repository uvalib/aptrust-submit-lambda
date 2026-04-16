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

	// if we are doing a bag purge, add the bag name to the suffix
	if be.EventName == uvaaptsbus.EventCommandBagPurge {
		pathPrefix = path.Join(pathPrefix, wf.BagId)
	}

	// get a complete list of all the files included in the specified path
	keys, err := s3Client.s3List(cfg.AssetBucket, pathPrefix)
	if err != nil {
		fmt.Printf("WARNING: listing S3 assets (%s), continuing\n", err.Error())
	}

	fmt.Printf("INFO: located %d assets in [s3://%s/%s]\n", len(keys), cfg.AssetBucket, pathPrefix)

	efsDir := path.Join(cfg.AssetFilesystem, pathPrefix)
	contents, err := os.ReadDir(efsDir)
	if err != nil {
		fmt.Printf("WARNING: listing cache assets (%s), continuing\n", err.Error())
	}

	fmt.Printf("INFO: located %d assets in cache [%s]\n", len(contents), efsDir)

	// purge the S3 assets
	err = purgeS3Assets(s3Client, cfg.AssetBucket, keys)
	if err != nil {
		fmt.Printf("WARNING: purging S3 assets (%s), continuing\n", err.Error())
	}

	// purge the cache
	err = purgeCacheAssets(efsDir, contents)
	if err != nil {
		fmt.Printf("WARNING: purging cache assets (%s), continuing\n", err.Error())
	}

	return nil
}

//
// end of file
//
