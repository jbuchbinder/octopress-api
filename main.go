package main

import (
	"flag"
	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const (
	DEFAULT_MIME_TYPE = "application/json"
)

var (
	bind       = flag.String("bind", ":8888", "Port/IP for binding interface")
	username   = flag.String("username", "admin", "Username for BASIC auth")
	password   = flag.String("password", "password", "Password for BASIC auth")
	gitcmd     = flag.String("git", "git", "Executable for git command")
	rakecmd    = flag.String("rake", "rake", "Executable for rake command")
	retmime    = flag.String("mime", DEFAULT_MIME_TYPE, "MIME type for JSON responses")
	uiLocation = flag.String("ui", "./ui", "Location of UI")
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

	log.Print("Octopress API v" + Version + " / API v" + CurrentApiVersion)

	// Compile all instances into sites
	log.Print("Compiling instances into sites")
	for i := range instances {
		log.Print("Processing site " + instances[i])
		site, err := getSite(instances[i])
		if err != nil {
			panic(err)
		}
		log.Print("Identified site name " + site.Name)
		MySitesMap[site.Key] = site
	}

	r := mux.NewRouter()

	// Handle UI prefix
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir(*uiLocation))))

	api := r.PathPrefix("/api").Subrouter()

	// Common to all version
	api.HandleFunc("/version", versionHandler).Methods("GET")

	// API Version 1.0
	subV1_0 := api.PathPrefix("/1.0").Subrouter()
	subV1_0.HandleFunc("/version", versionHandler).Methods("GET")
	subV1_0.HandleFunc("/sites", sitesHandler).Methods("GET")
	subV1_0.HandleFunc("/site/commit/{site}", gitCommitHandler).Methods("GET")
	subV1_0.HandleFunc("/site/deploy/{site}", deployHandler).Methods("GET")
	subV1_0.HandleFunc("/post/categories/{site}", postCategoriesHandler).Methods("GET")
	subV1_0.HandleFunc("/post/list/{site}", listPostsHandler).Methods("GET")
	subV1_0.HandleFunc("/post/new/{site}/{postname}", newPostHandler).Methods("GET")
	subV1_0.HandleFunc("/post/update/{site}/{slug}", updatePostHandler).Methods("POST")

	// Redirection for home
	r.HandleFunc("/", homeRedirectHandler)

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
