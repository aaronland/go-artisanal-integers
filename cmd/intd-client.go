package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/client"
	"log"
	"net/url"
)

func main() {

	var proto = flag.String("protocol", "http", "...")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")
	var path = flag.String("path", "/", "The port number to listen for requests on")

	flag.Parse()

	u := new(url.URL)

	u.Scheme = *proto
	u.Host = fmt.Sprintf("%s:%d", *host, *port)
	u.Path = *path

	_, err := url.Parse(u.String())

	if err != nil {
		log.Fatal(err)
	}

	cl, err := client.NewArtisanalClient(*proto, u)

	if err != nil {
		log.Fatal(err)
	}

	i, err := cl.NextInt()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(i)
}
