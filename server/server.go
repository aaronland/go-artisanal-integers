package server

import (
	"context"
	"fmt"
	aa_http "github.com/aaronland/go-artisanal-integers/http"
	"github.com/aaronland/go-artisanal-integers/database"
	aa_server "github.com/aaronland/go-http-server"
	_ "log"
	"net/http"
	"net/url"
	"strings"
)

type ArtisanalServer struct {
	aa_server.Server
	server  aa_server.Server
	database database.Database
	url     *url.URL
}

func NewArtisanalServer(ctx context.Context, uri string) (aa_server.Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	database_uri := q.Get("database")
	q.Del("database")
	
	if database_uri == "" {
		return nil, fmt.Errorf("Missing ?database= parameter")
	}

	db, err := database.NewDatabase(ctx, database_uri)

	if err != nil {
		return nil, err
	}

	u.RawQuery = q.Encode()
	uri = u.String()
	
	aa_svr, err := aa_server.NewServer(ctx, uri)

	if err != nil {
		return nil, err
	}

	svr := &ArtisanalServer{
		server:  aa_svr,
		database: db,
		url:     u,
	}

	return svr, nil
}

func (svr *ArtisanalServer) Address() string {
	return svr.server.Address()
}

func (svr *ArtisanalServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	integer_handler, err := aa_http.IntegerHandler(svr.database)

	if err != nil {
		return err
	}

	integer_path := svr.url.Path

	if !strings.HasPrefix(integer_path, "/") {
		integer_path = fmt.Sprintf("/%s", integer_path)
	}

	mux.(*http.ServeMux).Handle(integer_path, integer_handler)

	return svr.server.ListenAndServe(ctx, mux)
}
