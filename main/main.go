package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	urlshort "github.com/abdulkaderm36/gophercises/url-short"
	"github.com/abdulkaderm36/gophercises/url-short/main/db"
)

func main() {
	mux := defaultMux()

	ymlPath := flag.String("y", "", "yaml file")
	flag.Parse()

	database := db.InitDB()
	defer database.Close()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/abdulkader":     "https://github.com/abdulkaderm36",
	}
	database.InitData(pathsToUrls)
	// mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	dbHandler := urlshort.DBHandler(*database, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution
    `
	if *ymlPath != "" {
		data, err := os.ReadFile(*ymlPath)
		if err != nil {
			panic(err)
		}
		yaml = string(data)
	}
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), dbHandler)
	if err != nil {
		panic(err)
	}

	json := `
    [
        {
            "path": "/google",
            "url": "https://google.com"
        } 
    ]
    `

	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		database.Close()
		os.Exit(1)
	}()

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
