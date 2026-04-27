module github.com/uvalib/apt-submit-submission-status

go 1.26.0

// for local development
//replace github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao => ../../aptrust-submit-db-dao/uvaaptsdao

require (
	github.com/aws/aws-lambda-go v1.54.0
	github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao v0.0.0-20260427140620-2e49452502e3
)

require github.com/lib/pq v1.12.3 // indirect
