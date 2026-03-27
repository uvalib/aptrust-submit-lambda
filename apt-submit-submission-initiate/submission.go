//
//
//

package main

import (
	"slices"
	"sort"
	"strings"
)

// from the list of files included in the submission, find the manifests to determine which bags
// are included
func findIncludedBags(prefix string, suppliedFiles []string) []string {
	bags := make([]string, 0)

	for _, fname := range suppliedFiles {

		// is this a manifest
		if strings.HasSuffix(fname, manifestName) {
			// strip the prefix
			s := strings.TrimPrefix(fname, prefix)

			// split with the seperator
			bits := strings.Split(s, "/")
			if len(bits) == 3 {
				// save the bag name
				bags = append(bags, bits[1])
			}
		}
	}
	return bags
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
