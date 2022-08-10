package server

import (
	"context"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/http/api"
	"github.com/aaronland/go-artisanal-integers/service"
	aa_server "github.com/aaronland/go-http-server"
	_ "log"
	"net/http"
	"net/url"
	"strings"
)

type ArtisanalServer struct {
	aa_server.Server
	server  aa_server.Server
	service service.Service
	url     *url.URL
}

func NewArtisanalServer(ctx context.Context, uri string) (aa_server.Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	service_uri := q.Get("service")
	q.Del("service")

	if service_uri == "" {
		return nil, fmt.Errorf("Missing ?service= parameter")
	}

	svc, err := service.NewService(ctx, service_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new integer service, %w", err)
	}

	u.RawQuery = q.Encode()
	uri = u.String()

	aa_svr, err := aa_server.NewServer(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new server, %w", err)
	}

	svr := &ArtisanalServer{
		server:  aa_svr,
		service: svc,
		url:     u,
	}

	return svr, nil
}

func (svr *ArtisanalServer) Address() string {
	return svr.server.Address()
}

func (svr *ArtisanalServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	integer_handler, err := api.IntegerHandler(svr.service)

	if err != nil {
		return fmt.Errorf("Failed to create integer handler, %w", err)
	}

	integer_path := svr.url.Path

	if !strings.HasPrefix(integer_path, "/") {
		integer_path = fmt.Sprintf("/%s", integer_path)
	}

	mux.(*http.ServeMux).Handle(integer_path, integer_handler)

	return svr.server.ListenAndServe(ctx, mux)
}
