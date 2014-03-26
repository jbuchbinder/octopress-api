package main

import (
	"errors"
	"os"
	"regexp"
)

// fileExists is a convenience function to make it simpler to determine if a
// specified file exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// postExists determines whether or not a particular post is already present
// by checking the filename it should occupy.
func postExists(location, slug string) bool {
	return fileExists(location + "/source/_posts/" + slug + ".md")
}

// postSlugFromFilename extracts a "post slug" from a relative filename (i.e.
// one that has the form source/_posts/YYYY-MM-DD-TITLE-GOES-HERE.md)
func postSlugFromFilename(filename string) (string, error) {
	re, err := regexp.Compile(`source/_posts/(.+)\.md`)
	if err != nil {
		return "", err
	}
	res := re.FindStringSubmatch(filename)

	if len(res) < 2 {
		return "", errors.New("Unable to resolve post slug")
	}

	return res[1], nil
}
