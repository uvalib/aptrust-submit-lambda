package main

import (
	"fmt"
	"strings"
)

// Config defines all the service configuration parameters
type Config struct {
	// event bus definitions
	//BusName string // the message bus name

	// configuration needed for mail content
	ValidationFailedUrl     string // validation failure report URL
	ReconciliationFailedUrl string // reconciliation failure report URL
	ApprovalUrl             string // approval form URL

	// mailer configuration
	EmailSender    string   // the email sender
	SendEmail      bool     // do we send or just log
	DebugRecipient string   // the debug recipient if we want to send
	AdminEmail     string   // the administrator email
	MailCC         []string // always copy

	// SMTP configuration
	SMTPHost string // SMTP hostname
	SMTPPort int    // SMTP port number
	SMTPUser string // SMTP username
	SMTPPass string // SMTP password

	// database configuration
	DbHost     string // database host
	DbPort     int    // database port
	DbName     string // database name
	DbUser     string // database user
	DbPassword string // database password
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error

	// mail content configuration
	cfg.ApprovalUrl, err = ensureSetAndNonEmpty("APPROVAL_URL")
	if err != nil {
		return nil, err
	}
	cfg.ValidationFailedUrl, err = ensureSetAndNonEmpty("VALIDATION_FAIL_URL")
	if err != nil {
		return nil, err
	}
	cfg.ReconciliationFailedUrl, err = ensureSetAndNonEmpty("RECONCILIATION_FAIL_URL")
	if err != nil {
		return nil, err
	}

	// SMTP configuration
	cfg.SMTPHost, err = ensureSetAndNonEmpty("SMTP_HOST")
	if err != nil {
		return nil, err
	}
	cfg.SMTPPort, err = envToInt("SMTP_PORT")
	if err != nil {
		return nil, err
	}
	cfg.SMTPUser = envWithDefault("SMTP_USER", "")
	cfg.SMTPPass = envWithDefault("SMTP_PASSWORD", "")

	cfg.EmailSender, err = ensureSetAndNonEmpty("EMAIL_SENDER")
	if err != nil {
		return nil, err
	}
	cfg.AdminEmail, err = ensureSetAndNonEmpty("ADMIN_EMAIL")
	if err != nil {
		return nil, err
	}
	cfg.SendEmail, err = envToBool("EMAIL_SEND")
	if err != nil {
		return nil, err
	}
	cfg.DebugRecipient = envWithDefault("DEBUG_RECIPIENT", "")
	cfg.MailCC = strings.Split(envWithDefault("MAIL_CC", ""), ",")

	//cfg.BusName = envWithDefault("MESSAGE_BUS", "")

	// database configuration
	cfg.DbHost, err = ensureSetAndNonEmpty("DB_HOST")
	if err != nil {
		return nil, err
	}
	cfg.DbPort, err = envToInt("DB_PORT")
	if err != nil {
		return nil, err
	}
	cfg.DbName, err = ensureSetAndNonEmpty("DB_NAME")
	if err != nil {
		return nil, err
	}
	cfg.DbUser, err = ensureSetAndNonEmpty("DB_USER")
	if err != nil {
		return nil, err
	}
	cfg.DbPassword, err = ensureSetAndNonEmpty("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	// mail content configuration
	fmt.Printf("[CONFIG] ValidationFailedUrl     = [%s]\n", cfg.ValidationFailedUrl)
	fmt.Printf("[CONFIG] ReconciliationFailedUrl = [%s]\n", cfg.ReconciliationFailedUrl)
	fmt.Printf("[CONFIG] ApprovalUrl             = [%s]\n", cfg.ApprovalUrl)

	fmt.Printf("[CONFIG] EmailSender             = [%s]\n", cfg.EmailSender)
	fmt.Printf("[CONFIG] SendEmail               = [%t]\n", cfg.SendEmail)
	fmt.Printf("[CONFIG] DebugRecipient          = [%s]\n", cfg.DebugRecipient)
	fmt.Printf("[CONFIG] AdminEmail              = [%s]\n", cfg.AdminEmail)
	fmt.Printf("[CONFIG] MailCC                  = [%v]\n", cfg.MailCC)

	// SMTP configuration
	fmt.Printf("[CONFIG] SMTPHost                = [%s]\n", cfg.SMTPHost)
	fmt.Printf("[CONFIG] SMTPPort                = [%d]\n", cfg.SMTPPort)
	fmt.Printf("[CONFIG] SMTPUser                = [%s]\n", cfg.SMTPUser)
	fmt.Printf("[CONFIG] SMTPPass                = [%s]\n", cfg.SMTPPass)

	// database configuration
	fmt.Printf("[CONFIG] DbHost                  = [%s]\n", cfg.DbHost)
	fmt.Printf("[CONFIG] DbPort                  = [%d]\n", cfg.DbPort)
	fmt.Printf("[CONFIG] DbName                  = [%s]\n", cfg.DbName)
	fmt.Printf("[CONFIG] DbUser                  = [%s]\n", cfg.DbUser)
	fmt.Printf("[CONFIG] DbPassword              = [REDACTED]\n")

	return &cfg, nil
}

//
// end of file
//
