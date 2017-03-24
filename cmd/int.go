package main

import (
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"github.com/thisisaaronland/go-artisanal-integers/engine"
	"log"
	"os"
)

func main() {

	var db = flag.String("database", "", "...")
	var dsn = flag.String("dsn", "", "...")
	var last = flag.Int("last-id", 0, "...")

	flag.Parse()

	var eng artisanalinteger.Engine
	var err error

	switch *db {

	case "redis":
		eng, err = engine.NewRedisEngine(*dsn)
	case "summitdb":
		eng, err = engine.NewSummitDBEngine(*dsn)
	case "mysql":
		eng, err = engine.NewMySQLEngine(*dsn)
	default:
		log.Fatal("Invalid engine")
	}

	if err != nil {
		log.Fatal(err)
	}

	if *last != 0 {

		err = eng.SetLastId(int64(*last))

		if err != nil {
			log.Fatal(err)
		}
	}

	next, err := eng.NextId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(next)
	os.Exit(0)
}
