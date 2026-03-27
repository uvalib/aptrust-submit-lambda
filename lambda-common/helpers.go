//
//
//

package main

import (
	"github.com/aws/aws-lambda-go/events"
)

func apiGatewayProxyErrorResponse(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: status}, nil
}

//
// end of file
//
