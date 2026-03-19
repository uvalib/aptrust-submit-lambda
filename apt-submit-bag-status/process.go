//
// main message processing
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// create the data access object
	dao, err := uvaaptsdao.NewDao(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	if err != nil {
		return err
	}

	// cleanup on exit
	defer dao.Close()

	// get the bags
	bags, err := dao.GetBagsByStatus(BagStatusPendingIngest)
	if err != nil {
		if errors.As(err, &ErrBagNotFound) {
			fmt.Printf("INFO: no bags in '%s' status\n", BagStatusPendingIngest)
			return nil
		}
		return err
	}

	if len(bags) != 0 {
		httpClient := newHttpClient(1, cfg.HttpTimeout)
		// important, cleanup properly
		defer httpClient.CloseIdleConnections()

		for _, bg := range bags {

			//status := getAptStatus(httpClient, bg)
			//switch status {
			//case BagStatusPendingIngest:
			//default:
			//}
			//if status != BagStatusPendingIngest {
			//}

			if len(bg.ETag) != 0 {
				url := strings.Replace(cfg.APTStatusUrl, "{:etag}", bg.ETag, 1)

				fmt.Printf("DEBUG: checking for status of <%s/%s>\n", bg.Submission, bg.Name)
				res, err := httpGet(httpClient, url, cfg.APTUser, cfg.APTKey)
				if err == nil {
					//fmt.Printf("DEBUG: request [%s]\n", request.Body)
					r := AptStatusResponse{}
					err := json.Unmarshal(res, &r)
					if err != nil {
						fmt.Printf("ERROR: json.Unmarshal() failed (%s)\n", err.Error())
					} else {
						fmt.Printf("DEBUG: received [%v]\n", r)
					}
				} else {
					fmt.Printf("WARNING: getting bag status for <%s/%s> (%s)\n", bg.Submission, bg.Name, err.Error())
				}
			} else {
				fmt.Printf("ERROR: bag <%s/%s> has an empty etag, cannot check for status\n", bg.Submission, bg.Name)
			}
		}
	}

	return nil
}

//
// end of file
//
