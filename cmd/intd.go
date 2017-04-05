package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/thisisaaronland/go-artisanal-integers/service"
	"github.com/thisisaaronland/go-artisanal-integers/util"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

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

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		next, err := svc.NextInt()

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_next := strconv.FormatInt(next, 10)

		b := []byte(str_next)

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(b)))

		rsp.Write(b)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("listening on ", endpoint)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err = gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
