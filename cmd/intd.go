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
	var db = flag.String("engine", "", "...")
	var dsn = flag.String("dsn", "", "...")

	flag.Parse()

	eng, err := util.NewArtisanalEngine(*db, *dsn)

	if err != nil {
		log.Fatal(err)
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
