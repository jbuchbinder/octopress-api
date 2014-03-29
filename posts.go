package main

import (
	"errors"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"strings"
)

type PostObject struct {
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Filename   string   `json:"filename"`
	Permalink  string   `json:"permalink"`
	Author     string   `json:"author"`
	Date       string   `json:"date"`
	Categories []string `json:"categories"`
	FullText   string   `json:"fulltext"`
	YamlHeader string   `json:"header"`
	Body       string   `json:"text"`
}

func parsePost(site Site, slug string) (PostObject, error) {
	obj := PostObject{}

	fullPost := site.Location + "/source/_posts/" + slug + ".md"
	if !fileExists(fullPost) {
		return obj, errors.New(fullPost + " does not exist")
	}

	newpost, err := ioutil.ReadFile(fullPost)
	if err != nil {
		return obj, err
	}

	// Parse out the YAML
	postconfig := postYaml{}
	err = yaml.Unmarshal([]byte(newpost), &postconfig)
	if err != nil {
		return obj, err
	}

	// Bring over fields from YAML
	obj.Title = postconfig.Title
	obj.Permalink = postconfig.Permalink
	obj.Categories = postconfig.Categories
	if postconfig.Date != "" {
		obj.Date = postconfig.Date
	}

	obj.Filename = "source/_posts/" + slug + ".md"
	obj.FullText = string(newpost)

	// Figure out date from slug
	if obj.Date == "" {
		slugParts := strings.Split(slug, "-")
		if len(slugParts) > 3 {
			obj.Date = slugParts[0] + "-" + slugParts[1] + "-" + slugParts[2]
		}
	}

	// Fancy splitting routines here
	segments := strings.Split(string(newpost), "---\n")
	if len(segments) != 3 {
		return obj, errors.New("Unable to parse post " + slug)
	}
	obj.YamlHeader = segments[1]
	obj.Body = segments[2]

	return obj, nil
}
