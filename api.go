package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
