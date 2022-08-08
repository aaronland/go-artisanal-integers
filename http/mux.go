package http

import (
	"fmt"
	"github.com/aaronland/go-artisanal-integers/service"
	"github.com/aaronland/go-http-ping/v2"
	gohttp "net/http"
	"net/url"
	"strings"
)

func NewServeMux(s service.Service, u *url.URL) (*gohttp.ServeMux, error) {

	integer_handler, err := IntegerHandler(s)

	if err != nil {
		return nil, err
	}

	integer_path := u.Path

	if !strings.HasPrefix(integer_path, "/") {
		integer_path = fmt.Sprintf("/%s", integer_path)
	}

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return nil, err
	}

	ping_path := "/ping"

	mux := gohttp.NewServeMux()

	mux.Handle(integer_path, integer_handler)
	mux.Handle(ping_path, ping_handler)

	return mux, nil
}
