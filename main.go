package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
)

const shutdownTimeout = 5 * time.Second

var (
	port     = flag.Int("port", 8080, "Server port number")
	basePath = flag.String("basePath", "", "Base path for routing")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	log.Print("Starting...")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			BasePath string
			Name     string
		}{
			BasePath: *basePath,
			Name:     generate(),
		}
		w.WriteHeader(http.StatusOK)
		tpl.Execute(w, data) // nolint: errcheck,gosec
	}
	http.HandleFunc("/", indexHandler)

	fileHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileHandler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.DefaultServeMux,
	}

	// Wait for stop signal
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Listening at %s...", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed with error: %v", err)
	}
	log.Print("Shutdown")
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
