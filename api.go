package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// apiFail triggers a failure response from the API with a specified
// message.
func apiFail(w http.ResponseWriter, r *http.Request, message string) {
	resp := apiFailureResponse{
		Success: false,
		Message: message,
	}
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
	w.WriteHeader(http.StatusNoContent)
}

// Define all generic callback functions for mux router here

func sitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", *retmime)
	b, _ := json.Marshal(MySitesMap)
	fmt.Fprint(w, string(b))
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", *retmime)
	vars := mux.Vars(r)
	instance := vars["site"]

	resp := CmdResponse{}

	site, found := MySitesMap[instance]
	if !found {
		apiFail(w, r, "Unable to locate site '"+instance+"'")
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
		apiFail(w, r, err.Error())
		return
	}

	resp.Success = true
	resp.Message = out
	b, _ := json.Marshal(resp)
	fmt.Fprint(w, string(b))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	versionMap := map[string]string{
		"version":           Version,
		"currentApiVersion": CurrentApiVersion,
	}

	w.Header().Add("Content-Type", *retmime)
	b, _ := json.Marshal(versionMap)
	fmt.Fprint(w, string(b))
}
