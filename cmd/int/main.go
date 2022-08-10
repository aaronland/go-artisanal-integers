package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/aaronland/go-artisanal-integers/service"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	var service_uri = flag.String("service-uri", "memory://", "The data source name (dsn) for connecting to the artisanal integer service.")
	var continuous = flag.Bool("continuous", false, "Continuously mint integers. This is mostly only useful for debugging.")

	flag.Parse()

	ctx := context.Background()

	s, err := service.NewService(ctx, *service_uri)

	if err != nil {
		log.Fatal(err)
	}

	writers := []io.Writer{
		os.Stdout,
	}

	multi := io.MultiWriter(writers...)
	writer := bufio.NewWriter(multi)

	for {

		next, err := s.NextInt(ctx)

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
