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
	var bagList string

	flag.StringVar(&messageId, "messageid", "0-0-0-0", "Message identifier")
	flag.StringVar(&clientId, "cid", "", "The client identifier")
	flag.StringVar(&submissionId, "sid", "", "The submission identifier")
	flag.StringVar(&bagList, "baglist", "", "A file containing list of bags (json format)")
	flag.Parse()

	if len(clientId) == 0 || len(submissionId) == 0 || len(bagList) == 0 {
		fmt.Printf("ERROR: incorrect commandline, use --help for details\n")
		os.Exit(1)
	}

	req := events.APIGatewayProxyRequest{}
	req.QueryStringParameters = map[string]string{}

	bags, err := os.ReadFile(bagList)
	if err != nil {
		fmt.Printf("ERROR: could not read bags file %s (%s)\n", bagList, err.Error())
		os.Exit(1)
	}

	req.Body = fmt.Sprintf("{\"cid\":\"%s\",\"sid\":\"%s\",\"bag_folders\":%s}", clientId, submissionId, string(bags))

	resp, err := process(messageId, "api.gateway", req)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		//os.Exit(1)
	}

	fmt.Printf("INFO: response: %s\n", resp.Body)
	fmt.Printf("INFO: terminating with HTTP %d\n", resp.StatusCode)
}

//
// end of file
//
