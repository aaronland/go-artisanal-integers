package server

import (
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-artisanal-integers/http"
	"log"
	gohttp "net/http"
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

	handler, err := http.IntegerHandler(service)

	if err != nil {
		return err
	}

	log.Println("listening on", s.address)

	mux := gohttp.NewServeMux()
	mux.Handle("/", handler)

	return gohttp.ListenAndServe(s.address, mux)
}
