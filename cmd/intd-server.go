package main

import (
	"flag"
	"github.com/aaronland/go-artisanal-integers/application"
	"github.com/aaronland/go-artisanal-integers/engine"
	"log"
	"os"
)

func main() {

	var dsn = flag.String("dsn", "", "The data source name (dsn) for connecting to the artisanal integer engine.")

	flag.Parse()

	eng, err := engine.NewMemoryEngine(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = application.ServerApplication(eng)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
