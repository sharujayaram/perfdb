package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"bitbucket.org/tebeka/nrsc"
	"github.com/alexcesaro/log/stdlog"
)

var logger = stdlog.GetFromFlags()

var storage storageHandler

func requestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof("%s %s", r.Method, r.URL)
		handler.ServeHTTP(rw, r)
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Storage initialization
	var err error
	if storage, err = newMongoHandler(); err != nil {
		os.Exit(1)
	}

	// Static assets
	nrsc.Handle("/static/")

	// RESTful API
	http.Handle("/", newRouter())

	// Banner and launcher
	fmt.Println("\n\t:-:-: perfkeeper :-:-:\t\t\tserving http://0.0.0.0:8080/\n")
	logger.Critical(http.ListenAndServe("0.0.0.0:8080", requestLog(http.DefaultServeMux)))
}
