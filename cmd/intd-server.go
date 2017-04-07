package main

import (
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers/service"
	"github.com/thisisaaronland/go-artisanal-integers/util"
	"log"
	"os"
)

func main() {

	var proto = flag.String("protocol", "http", "The protocol for the server to implement. Valid options are: http,tcp.")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var db = flag.String("engine", "", "The name of the artisanal integer engine to use.")
	var dsn = flag.String("dsn", "", "The data source name (dsn) for connecting to the artisanal integer engine.")
	var last = flag.Int("set-last-int", 0, "Set the last known integer.")
	var offset = flag.Int("set-offset", 0, "Set the offset used to mint integers.")
	var increment = flag.Int("set-increment", 0, "Set the increment used to mint integers.")

	flag.Parse()

	eng, err := util.NewArtisanalEngine(*db, *dsn)

	if err != nil {
		log.Fatal(err)
	}

	if *last != 0 {

		err = eng.SetLastInt(int64(*last))

		if err != nil {
			log.Fatal(err)
		}
	}

	if *increment != 0 {

		err = eng.SetIncrement(int64(*increment))

		if err != nil {
			log.Fatal(err)
		}
	}

	if *offset != 0 {

		err = eng.SetOffset(int64(*offset))

		if err != nil {
			log.Fatal(err)
		}
	}

	svc, err := service.NewExampleService(eng)

	if err != nil {
		log.Fatal(err)
	}

	address := fmt.Sprintf("%s:%d", *host, *port)

	s, err := util.NewArtisanalServer(*proto, address)

	if err != nil {
		log.Fatal(err)
	}

	err = s.ListenAndServe(svc)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
