//
//
//

// include this on a cmdline build only
//go:build cmdline

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

func main() {

	var messageId string
	var source string
	var eventName string
	var clientId string
	var submissionId string
	var bagId string
	//var detail string
	var eventTime string

	flag.StringVar(&messageId, "messageid", "0-0-0-0", "Message identifier")
	flag.StringVar(&source, "source", "the.source", "Message source")
	flag.StringVar(&eventName, "eventname", "", "Event name")
	flag.StringVar(&eventTime, "eventtime", "", "Time of the event (optional)")
	flag.StringVar(&clientId, "cid", "", "The event client identifier (optional)")
	flag.StringVar(&submissionId, "sid", "", "The event submission identifier (optional)")
	flag.StringVar(&bagId, "bid", "", "The event bag identifier (optional)")
	//flag.StringVar(&detail, "detail", "", "Event detail, usually json (optional)")
	flag.Parse()

	if len(eventName) == 0 {
		fmt.Printf("ERROR: incorrect commandline, use --help for details\n")
		os.Exit(1)
	}

	ev := uvaaptsbus.UvaBusEvent{}
	ev.EventName = eventName
	ev.EventTime = eventTime
	ev.ClientId = clientId

	if len(submissionId) != 0 || len(bagId) != 0 {
		wf := uvaaptsbus.UvaWorkflowEvent{SubmissionId: submissionId, BagId: bagId}
		b, err := wf.Serialize()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
		ev.Detail = b
	}

	//ev.SubmissionId = submissionId
	//ev.BagId = bagId
	//if len(detail) != 0 {
	//	ev.Detail = json.RawMessage(detail)
	//}

	pl, _ := ev.Serialize()
	err := process(messageId, source, pl)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("INFO: terminating normally\n")
}

//
// end of file
//
