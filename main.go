package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jfrazelle/leeroy/jenkins"
)

const (
	VERSION = "v0.1.0"
)

var (
	certFile   string
	keyFile    string
	port       string
	configFile string
	debug      bool
	version    bool

	config Config
)

type Config struct {
	Jenkins      jenkins.Client `json:"jenkins"`
	BuildCommits string         `json:"build_commits"`
	GHToken      string         `json:"github_token"`
	Builds       []Build        `json:"builds"`
}

type Build struct {
	Repo    string `json:"github_repo"`
	Job     string `json:"jenkins_job_name"`
	Context string `json:"context"`
}

func init() {
	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")
	flag.StringVar(&certFile, "cert", "", "path to ssl certificate")
	flag.StringVar(&keyFile, "key", "", "path to ssl key")
	flag.StringVar(&port, "port", "80", "port to use")
	flag.StringVar(&configFile, "config", "/etc/leeroy/config.json", "path to config file")
	flag.Parse()
}

func main() {
	// set log level
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if version {
		fmt.Println(VERSION)
		return
	}

	// read the config file
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Errorf("config file does not exist: %s", configFile)
		return
	}
	c, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf("could not read config file: %v", err)
		return
	}
	if err := json.Unmarshal(c, &config); err != nil {
		log.Errorf("error parsing config file as json: %v", err)
		return
	}

	// create mux server
	mux := http.NewServeMux()

	// jenkins notification endpoint
	mux.HandleFunc("/notifications/jenkins", jenkinsHandler)

	// github webhooks endpoint
	mux.HandleFunc("/notifications/github", githubHandler)

	// retry build endpoint
	mux.HandleFunc("/build/retry", retryBuildHandler)

	// custom build endpoint
	mux.HandleFunc("/build/custom", customBuildHandler)

	// set up the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %q", port)
	if certFile != "" && keyFile != "" {
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}