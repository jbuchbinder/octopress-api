package main

import (
	"errors"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"strings"
)

var (
	titleSanitizer = strings.NewReplacer("/", "_", " ", "_", "&", "_", ".", "_")
)

type SitesMap map[string]Site

type Site struct {
	Name     string
	Location string
}

type ConfigYaml struct {
	Url   string
	Title string
}

func GetSite(location string) (Site, error) {
	// Open the _config.yml file, and look for the site name
	if !Exists(location + "/_config.yml") {
		return Site{}, errors.New("No _config.yml file found for " + location + " instance")
	}
	filedata, err := ioutil.ReadFile(location + "/_config.yml")
	if err != nil {
		return Site{}, err
	}

	config := ConfigYaml{}
	err = yaml.Unmarshal([]byte(filedata), &config)
	if err != nil {
		return Site{}, err
	}

	site := Site{}
	site.Name = titleSanitizer.Replace(config.Title)
	site.Location = location
	return site, nil
}
