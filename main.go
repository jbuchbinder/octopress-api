package main

import (
	"flag"
	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var (
	bind     = flag.String("bind", ":8888", "Port/IP for binding interface")
	username = flag.String("username", "admin", "Username for BASIC auth")
	password = flag.String("password", "password", "Password for BASIC auth")
	//log, _     = syslog.New(syslog.LOG_DEBUG, "octopress-api")
	MySitesMap = SitesMap{}
)

func main() {
	flag.Parse()
	instances := flag.Args()

	if len(instances) < 1 {
		log.Print("No Octopress instances were specified.")
		flag.PrintDefaults()
		return
	}

	// Compile all instances into sites
	log.Print("Compiling instances into sites")
	for i := range instances {
		log.Print("Processing site " + instances[i])
		site, err := GetSite(instances[i])
		if err != nil {
			panic(err)
		}
		log.Print("Identified site name " + site.Name)
		MySitesMap[site.Name] = site
	}

	r := mux.NewRouter()

	sub := r.PathPrefix("/api").Subrouter()

	sub.HandleFunc("/sites", sitesHandler).Methods("GET")
	sub.HandleFunc("/deploy/{site}", deployHandler).Methods("GET")

	s := &http.Server{
		Addr:           *bind,
		ReadTimeout:    90 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Handle authentication
	a := auth.NewBasicAuthenticator("Octopress API", SimpleSecurityProvider(*username, *password))
	http.Handle("/", a.Wrap(func(w http.ResponseWriter, ar *auth.AuthenticatedRequest) {
		r.ServeHTTP(w, &ar.Request)
	}))

	// Run actual server
	log.Print("Starting server on " + *bind)
	log.Print(s.ListenAndServe().Error())
}
