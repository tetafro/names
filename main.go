package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
)

var port = flag.Int("port", 8080, "Server port number")

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Name string
		}{
			Name: generate(),
		}
		w.WriteHeader(http.StatusOK)
		tpl.Execute(w, data) // nolint: errcheck,gosec
	}
	http.HandleFunc("/", indexHandler)

	fileHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileHandler))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Listening at %s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed with error: %v", err)
	}
}

// generate generates an adjective and a noun with the same first letter.
func generate() string {
	for {
		adj := petname.Adjective()
		noun := petname.Name()
		if adj[0] == noun[0] {
			return adj + " " + noun
		}
	}
}
