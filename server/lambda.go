package server

import (
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-artisanal-integers/http"
	"github.com/whosonfirst/algnhsa"
	_ "log"
	gohttp "net/http"
)

type LambdaServer struct {
	artisanalinteger.Server
}

func NewLambdaServer() (*LambdaServer, error) {

	server := LambdaServer{}

	return &server, nil
}

func (s *LambdaServer) ListenAndServe(service artisanalinteger.Service) error {

	handler, err := http.IntegerHandler(service)

	if err != nil {
		return err
	}

	mux := gohttp.NewServeMux()
	mux.HandleFunc("/", handler)

	algnhsa.ListenAndServe(mux, nil)
	return nil
}
