package client

import (
	"errors"
	"github.com/aaronland/go-artisanal-integers"
	"strings"
)

func NewArtisanalClient(proto string, address string) (artisanalinteger.Client, error) {

	var cl artisanalinteger.Client
	var err error

	switch strings.ToUpper(proto) {

	case "HTTP":

		if address == "" {
			address = "localhost:8080"
		}

		cl, err = NewHTTPClient(address)

	case "TCP":

		cl, err = NewTCPClient(address)

	default:
		return nil, errors.New("Invalid client protocol")
	}

	if err != nil {
		return nil, err
	}

	return cl, nil
}
