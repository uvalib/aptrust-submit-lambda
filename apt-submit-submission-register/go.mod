module github.com/uvalib/apt-submit-submission-register

go 1.26.0

// for local development
//replace github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao => ../../aptrust-submit-db-dao/uvaaptsdao

require (
	github.com/aws/aws-lambda-go v1.53.0
	github.com/rs/xid v1.6.0
	github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao v0.0.0-20260316181608-f47e243baf5e
)

require github.com/lib/pq v1.11.2 // indirect
