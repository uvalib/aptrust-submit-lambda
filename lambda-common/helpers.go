//
//
//

package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func apiGatewayProxyErrorResponse(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: status}, err
}

func readFile(path string) ([]string, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

//
// end of file
//
