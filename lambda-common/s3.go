//
// simple module to get and set parameter values in the ssm
//

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type uvaS3Client struct {
	client *s3.Client
}

func newS3Client() (*uvaS3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	c := uvaS3Client{}
	c.client = s3.NewFromConfig(cfg)
	return &c, nil
}

func (c *uvaS3Client) s3List(bucket string, key string) ([]string, error) {

	//fmt.Printf("INFO: s3 list [%s/%s]\n", bucket, key)
	start := time.Now()

	// query parameters
	params := &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &key,
	}

	// create a paginator
	var limit int32 = 1000
	paginate := s3.NewListObjectsV2Paginator(c.client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = limit
	})

	// make the result set
	result := make([]string, 0)

	// iterate through the pages
	for paginate.HasMorePages() {

		// get the next page
		page, err := paginate.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, o := range page.Contents {
			//fmt.Printf("DEBUG: found [%s]\n", *o.Key)
			result = append(result, *o.Key)
		}
	}

	duration := time.Since(start)
	fmt.Printf("INFO: s3 list [%s/%s] complete in %0.2f seconds\n", bucket, key, duration.Seconds())
	return result, nil
}

func (s *uvaS3Client) s3Remove(bucket string, key string) error {

	fmt.Printf("INFO: deleting [%s/%s]", bucket, key)
	start := time.Now()

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	duration := time.Since(start)
	fmt.Printf("INFO: delete [%s/%s] complete in %0.2f seconds (%s)\n", bucket, key, duration.Seconds(), s.statusText(err))
	return err
}

func (s *uvaS3Client) statusText(err error) string {
	if err == nil {
		return "ok"
	}
	return "ERR"
}

//
// end of file
//
