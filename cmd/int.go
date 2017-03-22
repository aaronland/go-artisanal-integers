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

	var eng artisanalinteger.Engine
	var err error

	switch *db {

	case "redis":
		eng, err = engine.NewRedisEngine()
	case "summitdb":
		eng, err = engine.NewSummitDBEngine()
	case "mysql":
		eng, err = engine.NewMySQLEngine()
	default:
		log.Fatal("Invalid engine")
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(eng.Next())
}
