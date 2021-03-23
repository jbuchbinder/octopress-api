package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type listPostsResponse map[string]listPostItem
type listPostItem struct {
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Filename   string   `json:"filename"`
	Permalink  string   `json:"permalink"`
	Author     string   `json:"author"`
	Date       string   `json:"date"`
	Categories []string `json:"categories"`
}

type postYaml struct {
	Author     string
	Permalink  string
	Title      string
	Categories []string
	Date       string
}

func postCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", *retmime)
	vars := mux.Vars(r)
	instance := vars["site"]

	site, found := MySitesMap[instance]
	if !found {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := categoriesForSite(site)
	if err != nil {
		apiFail(w, r, err.Error())
		return
	}

	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

func listPostsHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(listPostsResponse)

	w.Header().Add("Content-Type", *retmime)
	vars := mux.Vars(r)
	instance := vars["site"]

	site, found := MySitesMap[instance]
	if !found {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Get directory listing for Location /source/_posts/
	fis, err := ioutil.ReadDir(site.Location + "/source/_posts/")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
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

		// First three parts of slug...
		slugParts := strings.Split(item.Slug, "-")
		item.Date = slugParts[0] + "-" + slugParts[1] + "-" + slugParts[2]

		item.Author = postconfig.Author
		item.Title = postconfig.Title
		item.Permalink = postconfig.Permalink
		item.Categories = postconfig.Categories

		resp[slug] = item
	}

	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

type newPostResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	PostFile   string   `json:"file"`
	PostText   string   `json:"post"`
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Date       string   `json:"date"`
	PostHeader string   `json:"header"`
	Permalink  string   `json:"permalink"`
	Categories []string `json:"categories"`
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	resp := newPostResponse{}

	w.Header().Add("Content-Type", *retmime)
	vars := mux.Vars(r)
	instance := vars["site"]
	slug := vars["slug"]

	site, found := MySitesMap[instance]
	if !found {
		apiFail(w, r, "Unable to locate site '"+instance+"'")
		return
	}

	fullPost := site.Location + "/source/_posts/" + slug + ".md"
	if !fileExists(fullPost) {
		apiFail(w, r, fullPost+" does not exist")
		return
	}

	post, err := parsePost(site, slug)
	if err != nil {
		apiFail(w, r, err.Error())
		return
	}

	resp.Success = true
	resp.Message = ""
	resp.Title = post.Title
	resp.PostFile = post.Filename
	resp.PostHeader = post.YamlHeader
	resp.PostText = post.Body
	resp.Slug = slug
	resp.Categories = post.Categories
	resp.Permalink = post.Permalink
	resp.Date = post.Date
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	resp := newPostResponse{}

	w.Header().Add("Content-Type", *retmime)
	vars := mux.Vars(r)
	instance := vars["site"]
	postName, err := url.QueryUnescape(vars["postname"])
	if err != nil {
		apiFail(w, r, "Unable to decode post name '"+vars["postname"]+"'")
		return
	}

	site, found := MySitesMap[instance]
	if !found {
		apiFail(w, r, "Unable to locate site '"+instance+"'")
		return
	}

	log.Print("Executing " + `new_post["` + strings.Replace(postName, `"`, `\\"`, -1) + `"]`)
	cmd := *rakecmd
	args := []string{
		`new_post["` + strings.Replace(postName, `"`, `\\"`, -1) + `"]`,
	}
	out, err := RunCmd(site.Location, cmd, args)
	fmt.Println(out)
	log.Print("Completed RunCmd")
	if err != nil {
		apiFail(w, r, err.Error())
		return
	}

	// Locate the new post information
	re, err := regexp.Compile(`Creating new post: (.+)`)
	res := re.FindStringSubmatch(out)

	if len(res) < 2 {
		apiFail(w, r, "Unable to fetch filename from regex")
		return
	}

	log.Print(site.Location + "/" + res[1])
	newpost, err := ioutil.ReadFile(site.Location + "/" + res[1])
	if err != nil {
		apiFail(w, r, err.Error())
		return
	}

	slug, err := postSlugFromFilename(res[1])
	if err != nil {
		apiFail(w, r, err.Error())
		return
	}

	resp.Success = true
	resp.Message = out
	resp.PostFile = res[1]
	resp.PostText = string(newpost)
	resp.Slug = slug
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	resp := CmdResponse{}

	// Decode post body
	postbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apiFail(w, r, "Unable to decode post body")
		return
	}

	w.Header().Add("Content-Type", *retmime)
	vars := mux.Vars(r)
	instance := vars["site"]
	slug := vars["slug"]

	site, found := MySitesMap[instance]
	if !found {
		apiFail(w, r, "Unable to locate site '"+instance+"'")
		return
	}

	// Check to make sure this exists already
	if !postExists(site.Location, slug) {
		apiFail(w, r, slug+" does not exist.")
		return
	}

	// Write to post file
	log.Print(site.Location + "/source/_posts/" + slug)
	err = ioutil.WriteFile(site.Location+"/source/_posts/"+slug, postbody, 0777)
	if err != nil {
		apiFail(w, r, err.Error())
		return
	}

	resp.Success = true
	resp.Message = slug + " successfully written."
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}
