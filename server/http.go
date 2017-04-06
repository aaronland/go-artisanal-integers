package server

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/thisisaaronland/go-artisanal-integers"
	"log"
	"net/http"
	"strconv"
)

type HTTPServer struct {
	artisanalinteger.Server
	address string
}

func NewHTTPServer(address string) (*HTTPServer, error) {

	server := HTTPServer{
		address: address,
	}

	return &server, nil
}

func (s *HTTPServer) ListenAndServe(service artisanalinteger.Service) error {

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		next, err := service.NextInt()

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

	log.Println("listening on", s.address)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := gracehttp.Serve(&http.Server{Addr: s.address, Handler: mux})
	return err
}
