package main

import (
	"flag"
	"github.com/thisisaaronland/go-artisanal-integers/engine"
	"log"
)

func main() {

	var engine = flag.String("engine", "", "...")

	switch *engine {

	case "redis":
		log.Println(*engine)
	case "summitdb":
		log.Println(*engine)
	case "mysql":
		log.Println(*engine)
	default:
		log.Fatal("Invalid engine")
	}
}
