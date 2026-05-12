module github.com/uvalib/apt-submit-submission-register

go 1.26.0

// for local development
//replace github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao => ../../aptrust-submit-db-dao/uvaaptsdao

require (
	github.com/aws/aws-lambda-go v1.54.0
	github.com/rs/xid v1.6.0
	github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao v0.0.0-20260512174704-6007e89eae40
)

require github.com/lib/pq v1.12.3 // indirect
