package main

import (
	//"encoding/json"
	//"fmt"
	//	"github.com/gorilla/mux"
	//	"io/ioutil"
	"net/http"
)

// Define all callback functions for mux router here

func sitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	//b, err := json.Marshal()
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprint(w, "false")
	//	return
	//}
	//fmt.Fprint(w, string(b))
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
}
