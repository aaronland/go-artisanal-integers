package main

import (
	"bufio"
	"flag"
	"github.com/thisisaaronland/go-artisanal-integers/util"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	var db = flag.String("engine", "", "...")
	var dsn = flag.String("dsn", "", "...")
	var last = flag.Int("last-id", 0, "...")
	var offset = flag.Int("offset", 0, "...")
	var increment = flag.Int("increment", 0, "...")
	var continuous = flag.Bool("continuous", false, "...")

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

	writers := []io.Writer{
		os.Stdout,
	}

	multi := io.MultiWriter(writers...)
	writer := bufio.NewWriter(multi)

	for {

		next, err := eng.NextInt()

		if err != nil {
			log.Fatal(err)
		}

		str_next := strconv.FormatInt(next, 10)
		writer.WriteString(str_next + "\n")
		writer.Flush()

		if !*continuous {
			break
		}

	}

	os.Exit(0)
}
