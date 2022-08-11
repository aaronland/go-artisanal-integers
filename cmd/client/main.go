package main

import (
	"context"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/client"
	"github.com/sfomuseum/go-flags/flagset"
	"log"
)

var client_uri string

func main() {

	fs := flagset.NewFlagSet("integer")

	fs.StringVar(&client_uri, "client-uri", "http://localhost:8080/", "")

	flagset.Parse(fs)

	ctx := context.Background()

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		log.Fatalf("Failed to create new client, %v", err)
	}

	i, err := cl.NextInt(ctx)

	if err != nil {
		log.Fatalf("Failed to get next integer, %v", err)
	}

	fmt.Printf("%d", i)
}
