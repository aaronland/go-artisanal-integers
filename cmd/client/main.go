package main

import (
	"context"
	"fmt"
	"log"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/aaronland/go-artisanal-integers/client"	
)

var client_uri string

func main() {

	fs := flagset.NewFlagSet("integer")

	fs.StringVar(&client_uri, "client-uri", "http://localhost:8080/", "")

	flagset.Parse(fs)

	ctx := context.Background()

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		log.Fatalf("Failed to create new client, %w", err)
	}

	i, err := cl.NextInt(ctx)

	if err != nil {
		log.Fatalf("Failed to get next integer, %w", err)
	}

	fmt.Printf("%d", i)
}
