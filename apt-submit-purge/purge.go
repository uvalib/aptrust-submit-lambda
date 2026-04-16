//
//
//

package main

import (
	"fmt"
	"os"
	"path"
)

func purgeS3Assets(s3Client *uvaS3Client, bucket string, keys []string) error {

	for _, k := range keys {
		s := fmt.Sprintf("s3://%s", path.Join(bucket, k))
		fmt.Printf("INFO: removing [%s]\n", s)
		err := s3Client.s3Remove(bucket, k)
		if err != nil {
			fmt.Printf("ERROR: removing [%s] (%s), continuing\n", s, err.Error())
		}
	}

	return nil
}

func purgeCacheAssets(dir string, contents []os.DirEntry) error {

	for _, de := range contents {
		fmt.Printf("INFO: removing [%s]\n", path.Join(dir, de.Name()))
		err := os.RemoveAll(path.Join(dir, de.Name()))
		if err != nil {
			fmt.Printf("ERROR: removing [%s] (%s), continuing\n", path.Join(dir, de.Name()), err.Error())
		}
	}

	return nil
}

//
// end of file
//
