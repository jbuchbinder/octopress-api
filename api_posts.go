package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type listPostsResponse map[string]listPostItem
type listPostItem struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Date      string `json:"data"`
	Filename  string `json:"filename"`
	Permalink string `json:"permalink"`
	Author    string `json:"author"`
}

type postYaml struct {
	Author    string
	Permalink string
	Title     string
}

func listPostsHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(listPostsResponse)

	w.Header().Add("Content-Type", "application/json")
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

		item.Author = postconfig.Author
		item.Title = postconfig.Title
		item.Permalink = postconfig.Permalink

		resp[slug] = item
	}

	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

type newPostResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	PostFile string `json:"file"`
	PostText string `json:"post"`
	Slug     string `json:"slug"`
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	resp := newPostResponse{}

	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	instance := vars["site"]
	postName, err := url.QueryUnescape(vars["postname"])
	if err != nil {
		resp.Success = false
		resp.Message = "Unable to decode post name '" + vars["postname"] + "'"
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	site, found := MySitesMap[instance]
	if !found {
		resp.Success = false
		resp.Message = "Unable to locate site '" + instance + "'"
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
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
		resp.Success = false
		resp.Message = err.Error()
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Locate the new post information
	re, err := regexp.Compile(`Creating new post: (.+)`)
	res := re.FindStringSubmatch(out)

	if len(res) < 2 {
		resp.Success = false
		resp.Message = "Unable to fetch filename from regex"
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Print(site.Location + "/" + res[1])
	newpost, err := ioutil.ReadFile(site.Location + "/" + res[1])
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	slug, err := postSlugFromFilename(res[1])
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
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
		resp.Success = false
		resp.Message = "Unable to decode post body"
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	instance := vars["site"]
	slug := vars["slug"]

	site, found := MySitesMap[instance]
	if !found {
		resp.Success = false
		resp.Message = "Unable to locate site '" + instance + "'"
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check to make sure this exists already
	if !postExists(site.Location, slug) {
		resp.Success = false
		resp.Message = slug + " does not exist."
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Write to post file
	log.Print(site.Location + "/source/_posts/" + slug)
	err = ioutil.WriteFile(site.Location+"/source/_posts/"+slug, postbody, 0777)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp.Success = true
	resp.Message = slug + " successfully written."
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}
