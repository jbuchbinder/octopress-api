package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Define all callback functions for mux router here

func gitCommitHandler(w http.ResponseWriter, r *http.Request) {
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

	// Two steps involved in this, both with commands.

	// 1) Make sure all sources/_posts/*.md files are properly versioned
	cmdAdd := *gitcmd
	argsAdd := []string{
		"add",
		"source/_posts/*.md",
	}
	// We honestly don't care if this works or not. Most of the time, it's not
	// even necessary.
	RunCmd(site.Location, cmdAdd, argsAdd)

	// 2) Issue git commit command

	cmdGit := *gitcmd
	argsGit := []string{
		"commit",
		"source/_posts",
	}
	outGit, err := RunCmd(site.Location, cmdGit, argsGit)
	log.Print("Completed RunCmd for git commit")
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		b, _ := json.Marshal(resp)
		fmt.Fprint(w, string(b))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp.Success = true
	resp.Message = outGit
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}
