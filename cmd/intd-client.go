package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/client"
	"log"
)

func main() {

	var proto = flag.String("protocol", "http", "...")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	address := fmt.Sprintf("%s:%d", *host, *port)

	cl, err := client.NewArtisanalClient(*proto, address)

	if err != nil {
		log.Fatal(err)
	}

	i, err := cl.NextInt()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(i)
}
