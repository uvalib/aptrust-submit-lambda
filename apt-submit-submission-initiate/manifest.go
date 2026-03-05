//
//
//

package main

import (
	"fmt"
	"strings"
)

type ManifestRow struct {
	hash string
	file string
	bag  string
}

func manifestContents(s3client *uvaS3Client, bucket string, prefix string, bagName string) ([]ManifestRow, error) {
	manifestKey := fmt.Sprintf("%s/%s/%s", prefix, bagName, manifestName)
	localName := fmt.Sprintf("%s/%s-%s", tempFilesystem, bagName, manifestName)

	// get the manifest
	err := s3client.s3Get(bucket, manifestKey, localName)
	if err != nil {
		return nil, err
	}

	lines, err := readFile(localName)
	if err != nil {
		return nil, err
	}

	//
	// manifests are a hash followed by two spaces followed by a filename (which could contain spaces)
	//

	results := make([]ManifestRow, 0)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		subs := strings.SplitN(line, " ", 2)
		if len(subs) == 2 {
			hash := strings.TrimSpace(subs[0])
			name := strings.TrimSpace(subs[1])
			ml := ManifestRow{hash: hash, file: name, bag: bagName}
			results = append(results, ml)
		}
	}

	return results, nil
}

//
// end of file
//
