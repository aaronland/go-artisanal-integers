package main

import (
	"bufio"
	"flag"
	"github.com/thisisaaronland/go-artisanal-integers"
	"github.com/thisisaaronland/go-artisanal-integers/engine"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	var db = flag.String("engine", "", "...")
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

	writers := []io.Writer{
		os.Stdout,
	}

	multi := io.MultiWriter(writers...)
	writer := bufio.NewWriter(multi)

	str_next := strconv.FormatInt(next, 10)
	writer.WriteString(str_next + "\n")
	writer.Flush()

	os.Exit(0)
}
