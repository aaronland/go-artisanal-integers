package main

import (
	"context"
	"github.com/aaronland/go-artisanal-integers/server"
	"github.com/sfomuseum/go-flags/flagset"
	"log"
)

var server_uri string

func main() {

	fs := flagset.NewFlagSet("integer")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080?service=memory://", "")

	flagset.Parse(fs)

	ctx := context.Background()

	s, err := server.NewArtisanalIntegerServer(ctx, server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listen on %s\n", s.Address())

	err = s.ListenAndServe(ctx)

	if err != nil {
		log.Fatalf("Failed to serve requests, %v", err)
	}
}
