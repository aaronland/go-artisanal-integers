package server

import (
	"errors"
	"github.com/aaronland/go-artisanal-integers"
	"strings"
)

func NewArtisanalServer(proto string, address string) (artisanalinteger.Server, error) {

	var svr artisanalinteger.Server
	var err error

	switch strings.ToUpper(proto) {

	case "HTTP":

		if address == "" {
			address = "localhost:8080"
		}

		svr, err = NewHTTPServer(address)

	case "LAMBDA":

		svr, err = NewLambdaServer()

	case "TCP":

		svr, err = NewTCPServer(address)

	default:
		return nil, errors.New("Invalid server protocol")
	}

	if err != nil {
		return nil, err
	}

	return svr, nil
}
