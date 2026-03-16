//
//
//

package main

import (
	"slices"
	"sort"
	"strings"

	"github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao"
)

// from the list of files included in the submission, find the discrete set of bags included
func findIncludedBags(prefix string, suppliedFiles []string) []string {
	results := make([]string, 0)

	for _, fname := range suppliedFiles {
		// strip the prefix
		s := strings.TrimPrefix(fname, prefix)

		// split with the seperator
		bits := strings.Split(s, "/")
		bagName := bits[1]
		if slices.Contains(results, bagName) == false {
			results = append(results, bagName)
		}
	}
	return results
}

// create the bags in the database
func createDBBags(dao *uvaaptsdao.Dao, bagList []string, sid string) error {

	// create the bags
	for _, bagName := range bagList {
		err := dao.AddBag(bagName, sid)
		if err != nil {
			return err
		}
	}
	return nil
}

// create the files in the database
func createDBFiles(dao *uvaaptsdao.Dao, fileList []ManifestRow, sid string) error {

	// create the files
	for _, mr := range fileList {
		err := dao.AddFile(mr.file, mr.hash, sid, mr.bag)
		if err != nil {
			return err
		}
	}
	return nil
}

// returns true of 2 lists of strings are identical
func areIdentical(a, b []string) bool {

	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	return slices.Equal(a, b)
}

//
// end of file
//
