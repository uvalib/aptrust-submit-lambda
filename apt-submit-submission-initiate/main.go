//
//
//

// include this on a cmdline build only
//go:build cmdline
// +build cmdline

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func main() {

	var messageId string
	var clientId string
	var submissionId string

	flag.StringVar(&messageId, "messageid", "0-0-0-0", "Message identifier")
	flag.StringVar(&clientId, "cid", "", "The client identifier")
	flag.StringVar(&submissionId, "sid", "", "The submission identifier")
	flag.Parse()

	if len(clientId) == 0 || len(submissionId) == 0 {
		fmt.Printf("ERROR: incorrect commandline, use --help for details\n")
		os.Exit(1)
	}

	req := events.APIGatewayProxyRequest{}
	req.QueryStringParameters = map[string]string{}
	req.QueryStringParameters["cid"] = clientId
	req.QueryStringParameters["sid"] = submissionId

	//	req.Body = "{\"bag_folders\":[\"LibraETD-7d278v497\", \"LibraETD-n009w3930\", \"LibraETD-z890rv91v\"]}"
	req.Body = "{\"bag_folders\":[\"LibraETD-ks65hd74b\"]}"

	resp, err := process(messageId, "api.gateway", req)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("INFO: response: %s\n", resp.Body)
	fmt.Printf("INFO: terminating with HTTP %d\n", resp.StatusCode)
}

//
// end of file
//
