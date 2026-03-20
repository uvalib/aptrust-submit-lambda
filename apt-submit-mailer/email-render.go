//
//
//

package main

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

// templates holds our email templates
//
//go:embed templates/*
var templates embed.FS

func renderSubjectAndBody(cfg *Config, recipient string, be *uvaaptsbus.UvaBusEvent, wf *uvaaptsbus.UvaWorkflowEvent) (string, string, error) {

	var templateFile string
	var subject string
	var url string
	switch be.EventName {
	case uvaaptsbus.EventSubmissionApprove:
		templateFile = "templates/submission-approve.template"
		subject = "APTrust submission approval required"
		url = cfg.ApprovalUrl

	case uvaaptsbus.EventSubmissionReconcileFail:
		templateFile = "templates/submission-conflict.template"
		subject = "APTrust submission conflicts encountered; investigation required"
		url = cfg.ConflictUrl
	}

	// substitute the submission id into the URL
	url = strings.Replace(url, "{{:sid}}", wf.SubmissionId, -1)

	// read the template
	templateStr, err := templates.ReadFile(templateFile)
	if err != nil {
		return "", "", err
	}

	// parse the templateFile
	tmpl, err := template.New("email").Parse(string(templateStr))
	if err != nil {
		return "", "", err
	}

	// variables required by the templates
	type Attributes struct {
		Recipient  string // email recipient
		Submission string // submission identifier
		Url        string // appropriate management URL
		Sender     string // the sender
	}

	//	populate the attributes
	attribs := Attributes{
		Recipient:  recipient,
		Submission: wf.SubmissionId,
		Url:        url,
		Sender:     cfg.EmailSender,
	}

	// render the template
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return "", "", err
	}

	return subject, renderedBuffer.String(), nil
}

//
// end of file
//
