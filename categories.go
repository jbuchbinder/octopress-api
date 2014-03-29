package main

import (
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"sort"
	"strings"
)

func categoriesForSite(site Site) ([]string, error) {
	resp := make([]string, 0)
	posts := make(listPostsResponse)

	// Get directory listing for Location /source/_posts/
	fis, err := ioutil.ReadDir(site.Location + "/source/_posts/")
	if err != nil {
		return resp, err
	}
	for iter := 0; iter < len(fis); iter++ {
		slug := strings.Replace(fis[iter].Name(), ".md", "", 1)
		item := listPostItem{
			Slug:     slug,
			Filename: fis[iter].Name(),
		}

		// Pull rest of information from yaml
		filedata, err := ioutil.ReadFile(site.Location + "/source/_posts/" + fis[iter].Name())
		if err != nil {
			continue
		}
		postconfig := postYaml{}
		err = yaml.Unmarshal([]byte(filedata), &postconfig)
		if err != nil {
			continue
		}

		item.Author = postconfig.Author
		item.Title = postconfig.Title
		item.Permalink = postconfig.Permalink
		item.Categories = postconfig.Categories

		posts[slug] = item
	}

	cMap := make(map[string]bool)

	// Get unique categories as map keys
	for k := range posts {
		for cIdx := range posts[k].Categories {
			if !cMap[posts[k].Categories[cIdx]] {
				cMap[posts[k].Categories[cIdx]] = true
			}
		}
	}

	// Pull out keys into an array to return
	for ck, _ := range cMap {
		resp = append(resp, ck)
	}

	// But make sure we do an alpha sort first
	sort.Strings(resp)

	return resp, nil
}
