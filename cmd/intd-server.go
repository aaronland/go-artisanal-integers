package main

import (
	"github.com/aaronland/go-artisanal-integers/application"
	"github.com/aaronland/go-artisanal-integers/engine"
	"log"
	"os"
)

func main() {

	flags := application.NewServerApplicationFlags()

	/*
		var dsn string
		flags.StringVar(&dsn, "dsn", "example", "The data source name (dsn) for connecting to the artisanal integer engine.")

		application.ParseFlags(flags)

		eng, err := engine.NewMemoryEngine(dsn)
	*/

	eng, err := engine.NewMemoryEngine("")

	if err != nil {
		log.Fatal(err)
	}

	app, err := application.NewServerApplication(eng)

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(flags)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
