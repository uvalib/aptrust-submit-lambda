module github.com/uvalib/apt-submit-submission-register

go 1.26.0

require (
	github.com/aws/aws-lambda-go v1.53.0
	github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao v0.0.0-20260309182453-229fc621ff6b
)

require (
	github.com/lib/pq v1.11.2 // indirect
	github.com/rs/xid v1.6.0 // indirect
)

// for local development
//replace github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao => ../../aptrust-submit-db-dao/uvaaptsdao
