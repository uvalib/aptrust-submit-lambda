//
//
//

package main

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
)

// templates holds our email templates
//
//go:embed templates/*
var templates embed.FS

func renderSubjectAndBody(cfg *Config, clientName string, recipient string, be *uvaaptsbus.UvaBusEvent, wf *uvaaptsbus.UvaWorkflowEvent) (string, string, error) {

	client := strings.Title(clientName)
	var templateFile string
	var subject string
	var url string
	switch be.EventName {
	case uvaaptsbus.EventSubmissionApprove:
		templateFile = "templates/submission-approve.template"
		subject = fmt.Sprintf("Approval required for %s APTrust submission", client)
		url = cfg.ApprovalUrl

	case uvaaptsbus.EventSubmissionValidateFail:
		templateFile = "templates/submission-validate-fail.template"
		subject = fmt.Sprintf("Validation failures for %s APTrust submission; investigation is required", client)
		url = cfg.ValidationFailedUrl

	case uvaaptsbus.EventSubmissionReconcileFail:
		templateFile = "templates/submission-reconcile-fail.template"
		subject = fmt.Sprintf("Content conflicts identified for %s APTrust submission; investigation is required", client)
		url = cfg.ReconciliationFailedUrl
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
		Client     string // the submitting client (work type)
	}

	//	populate the attributes
	attribs := Attributes{
		Recipient:  recipient,
		Submission: wf.SubmissionId,
		Url:        url,
		Sender:     cfg.EmailSender,
		Client:     client,
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
