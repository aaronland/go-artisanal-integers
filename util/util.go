package util

import (
	"errors"
	"github.com/thisisaaronland/go-artisanal-integers"
	"github.com/thisisaaronland/go-artisanal-integers/client"
	"github.com/thisisaaronland/go-artisanal-integers/server"
)

func NewArtisanalServer(proto string, address string) (artisanalinteger.Server, error) {

	var svr artisanalinteger.Server
	var err error

	switch proto {

	case "http":

		if address == "" {
			address = "localhost:8080"
		}

		svr, err = server.NewHTTPServer(address)

	case "tcp":

		svr, err = server.NewTCPServer(address)

	default:
		return nil, errors.New("Invalid protocol")
	}

	if err != nil {
		return nil, err
	}

	return svr, nil

}

func NewArtisanalClient(proto string, address string) (artisanalinteger.Client, error) {

	var cl artisanalinteger.Client
	var err error

	switch proto {

	case "http":

		if address == "" {
			address = "localhost:8080"
		}

		cl, err = client.NewHTTPClient(address)

	case "tcp":

		cl, err = client.NewTCPClient(address)

	default:
		return nil, errors.New("Invalid protocol")
	}

	if err != nil {
		return nil, err
	}

	return cl, nil
}
