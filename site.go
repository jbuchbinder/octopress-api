package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	titleSanitizer = strings.NewReplacer("/", "_", " ", "_", "&", "_", ".", "_")
)

type CmdResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SitesMap map[string]Site

type Site struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Location string `json:"location"`
}

type ConfigYaml struct {
	Url   string
	Title string
}

func getSite(location string) (Site, error) {
	// Open the _config.yml file, and look for the site name
	if !fileExists(location + "/_config.yml") {
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
	site.Name = config.Title
	site.Key = titleSanitizer.Replace(config.Title)
	site.Location = location
	return site, nil
}
