package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Templates
var files = []string{"./web/template/counter.html"}
var templ = template.Must(template.ParseFiles(files...))

var count = 0

func handleIncrement(w http.ResponseWriter, req *http.Request) {
	count += 1
	io.WriteString(w, fmt.Sprintf("%d", count))
}

type CountState struct {
	Count int
}

func handleCounter(w http.ResponseWriter, req *http.Request) {
	countState := CountState{count}

	err := templ.ExecuteTemplate(w, "counter", countState)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("home:", err)
	}
}

func main() {
	// Serve up static files, for stylesheets, htmx
	fs := http.FileServer(http.Dir("./web/static/"))

	// Setup our routes
	route := mux.NewRouter()
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	route.HandleFunc("/increment", handleIncrement)
	route.HandleFunc("/counter", handleCounter)
	http.Handle("/", route)

	// Start HTTP server
	port := "4000"
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
