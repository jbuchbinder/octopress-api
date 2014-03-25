package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type CmdResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Define all callback functions for mux router here

func sitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	b, _ := json.Marshal(MySitesMap)
	fmt.Fprint(w, string(b))
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	instance := vars["site"]

	resp := CmdResponse{}

	site, found := MySitesMap[instance]
	if !found {
		resp.Success = false
		resp.Message = "Unable to locate site '" + instance + "'"
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cmd := *rakecmd
	args := []string{
		"gen_deploy",
		//"generate",
	}
	out, err := RunCmd(site.Location, cmd, args)
	log.Print("Completed RunCmd")
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
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

type newPostResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	PostFile string `json:"file"`
	PostText string `json:"post"`
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

	resp.Success = true
	resp.Message = out
	resp.PostFile = res[1]
	resp.PostText = string(newpost)
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	versionMap := map[string]string{
		"version":           Version,
		"currentApiVersion": CurrentApiVersion,
	}

	w.Header().Add("Content-Type", "application/json")
	b, _ := json.Marshal(versionMap)
	fmt.Fprint(w, string(b))
}
