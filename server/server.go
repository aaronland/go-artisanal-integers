package server

import (
	"context"
	aa_server "github.com/aaronland/go-http-server"
	_ "log"
	gohttp "net/http"
	gourl "net/url"
)

type ArtisanalServer struct {
	aa_server Server
	url       *gourl.URL
}

func NewArtisanalServer(ctx context.Context, uri string) (aa_server.Server, error) {

	u := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	service_uri := q.Get("service")

	if service_uri == "" {
		return nil, errors.New("Missing ?service= parameter")
	}

	svc, err := service.NewService(ctx, service_uri)

	if err != nil {
		return nil, err
	}

	aa_svr, err := aa_server.NewServer(ctx, uri)

	if err != nil {
		return nil, err
	}

	svr := &ArtisanalServer{
		server:  aa_svr,
		service: svc,
	}

	return &svr, nil
}

func (svr *ArtisanalServer) Address() string {
	return svr.server.Address()
}

func (svr *ArtisanalServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {

	integer_handler, err := http.IntegerHandler(s)

	if err != nil {
		return nil, err
	}

	integer_path := u.Path

	if !strings.HasPrefix(integer_path, "/") {
		integer_path = fmt.Sprintf("/%s", integer_path)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		return nil, err
	}

	ping_path := "/ping"

	mux.Handle(integer_path, integer_handler)
	mux.Handle(ping_path, ping_handler)

	return svr.server.ListenAndServe(ctx, mux)
}
