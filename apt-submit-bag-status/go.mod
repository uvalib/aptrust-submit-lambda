module github.com/uvalib/apt-submit-bag-status

go 1.26.0

// for local development
//replace github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao => ../../aptrust-submit-db-dao/uvaaptsdao

require (
	github.com/aws/aws-lambda-go v1.54.0
	github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao v0.0.0-20260326165616-55d8e9deb16a
)

require github.com/lib/pq v1.12.1 // indirect
