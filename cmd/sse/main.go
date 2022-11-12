package main

import (
	"github.com/alexzvon/job-sse/internal/handler"
	"github.com/alexzvon/job-sse/internal/splitter"
	"log"
	"net/http"
)

func main() {
	s := splitter.NewSseSplitter()
	h := handler.NewSseHendler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/listen", h.ListenSseHandler)
	mux.HandleFunc("/say", h.SaySseHandler)

	server := &http.Server{
		Addr:    ":8085",
		Handler: mux,
	}

	log.Println(server.ListenAndServe())
}
