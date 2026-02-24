//
//
//

package main

import (
	"database/sql"
	"fmt"
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

//
// end of file
//
