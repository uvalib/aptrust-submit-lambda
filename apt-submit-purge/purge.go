//
//
//

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// the client and submission identifiers arrive from the event bus and are used to build
// a filesystem path that is recursively deleted; only accept plain path segments
func ensureSafeIdentifier(name string, value string) error {

	if len(value) == 0 {
		err := fmt.Errorf("%s is empty", name)
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	if value == "." || value == ".." || strings.ContainsAny(value, `/\`) {
		err := fmt.Errorf("%s contains unexpected characters [%s]", name, value)
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	return nil
}

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
