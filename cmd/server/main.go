package main

import (
	"context"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/aaronland/go-artisanal-integers/server"
	"net/http"
	"log"
)

var server_uri string

func main() {

	fs := flagset.NewFlagSet("integer")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080?service=memory://", "")

	flagset.Parse(fs)

	ctx := context.Background()
	
	s, err := server.NewArtisanalServer(ctx, server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listen on %s\n", s.Address())
	
	mux := http.NewServeMux()

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to serve requests, %v", err)
	}
}
