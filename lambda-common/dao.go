//
//
//

package main

import (
	"database/sql"
	"fmt"
	"github.com/rs/xid"
	"log"

	// postgres
	_ "github.com/lib/pq"
)

type Dao struct {
	log     *log.Logger // logger
	*sql.DB             // database connection
}

func newDao(cfg *Config) (*Dao, error) {

	// connection attributes
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	// connect and ensure success
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Printf("ERROR: unable to open database (%s)\n", err.Error())
		return nil, err
	}

	// try a ping before declaring victory
	if err = db.Ping(); err != nil {
		fmt.Printf("ERROR: unable to ping database (%s)\n", err.Error())
		return nil, err
	}

	// all good
	return &Dao{
		//log:             c.Log,
		DB: db,
	}, nil
}

// Check -- check our database health
func (dao *Dao) Check() error {
	return dao.Ping()
}

// GetSubmissionStatus -- get the status of the specified submission
func (dao *Dao) GetSubmissionStatus(sid string) (*SubmissionStatus, error) {

	rows, err := dao.Query("SELECT submission, status, updated_at FROM submission_status WHERE submission = $1 LIMIT 1", sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss, err := submissionStatusQueryResults(rows)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

// GetSubmission -- get the specified submission
func (dao *Dao) GetSubmission(sid string) (*Submission, error) {

	rows, err := dao.Query("SELECT identifier, client, created_at FROM submissions WHERE identifier = $1 LIMIT 1", sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s, err := submissionQueryResults(rows)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// GetClient -- get the client details for the specified identifier
func (dao *Dao) GetClient(cid string) (*Client, error) {

	rows, err := dao.Query("SELECT name, identifier, created_at FROM clients WHERE identifier = $1 LIMIT 1", cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	c, err := clientQueryResults(rows)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// CreateSubmission -- create a new submission for the specified client
func (dao *Dao) CreateSubmission(client string) (*Submission, error) {

	// insert into submissions
	stmt1, err := dao.Prepare("INSERT INTO submissions( identifier, client ) VALUES( $1,$2 )")
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	newIdentifier := newSubmissionIdentifier()
	err = execPrepared(stmt1, newIdentifier, client)
	if err != nil {
		return nil, err
	}

	// insert into submissions_status
	stmt2, err := dao.Prepare("INSERT INTO submission_status( submission, status ) VALUES( $1,$2 )")
	if err != nil {
		return nil, err
	}
	defer stmt2.Close()

	err = execPrepared(stmt2, newIdentifier, SubmissionStatusRegistered)
	if err != nil {
		return nil, err
	}

	// get the submission details
	s, err := dao.GetSubmission(newIdentifier)
	if err != nil {
		return nil, err
	}
	return s, nil
}

//
// internal helpers
//

func submissionQueryResults(rows *sql.Rows) (*Submission, error) {
	results := Submission{}
	count := 0

	for rows.Next() {
		err := rows.Scan(&results.Identifier, &results.Client, &results.Created)
		if err != nil {
			return nil, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// check for not found
	if count == 0 {
		return nil, fmt.Errorf("%q: %w", "object(s) not found", ErrSubmissionNotFound)
	}

	//logDebug(log, fmt.Sprintf("found %d object(s)", count))
	return &results, nil
}

func submissionStatusQueryResults(rows *sql.Rows) (*SubmissionStatus, error) {
	results := SubmissionStatus{}
	count := 0

	for rows.Next() {
		err := rows.Scan(&results.Identifier, &results.Status, &results.Updated)
		if err != nil {
			return nil, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// check for not found
	if count == 0 {
		return nil, fmt.Errorf("%q: %w", "object(s) not found", ErrSubmissionNotFound)
	}

	//logDebug(log, fmt.Sprintf("found %d object(s)", count))
	return &results, nil
}

func clientQueryResults(rows *sql.Rows) (*Client, error) {
	results := Client{}
	count := 0

	for rows.Next() {
		err := rows.Scan(&results.Name, &results.Identifier, &results.Created)
		if err != nil {
			return nil, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// check for not found
	if count == 0 {
		return nil, fmt.Errorf("%q: %w", "object(s) not found", ErrClientNotFound)
	}

	//logDebug(log, fmt.Sprintf("found %d object(s)", count))
	return &results, nil
}

func execPrepared(stmt *sql.Stmt, values ...any) error {
	_, err := stmt.Exec(values...)
	return err
}

func newSubmissionIdentifier() string {
	return fmt.Sprintf("sid-%s", xid.New().String())
}

//
// end of file
//
