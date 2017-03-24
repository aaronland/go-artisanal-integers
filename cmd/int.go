package main

import (
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"github.com/thisisaaronland/go-artisanal-integers/engine"
	"log"
)

func main() {

	var db = flag.String("database", "", "...")
	var dsn = flag.String("dsn", "", "...")
	var incr = flag.Int("auto-increment", 0, "...")

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

	if *incr != 0 {

		err = eng.Set(int64(*incr))

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println(eng.Next())
}
