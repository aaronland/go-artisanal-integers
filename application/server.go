package application

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/engine"
	"github.com/aaronland/go-artisanal-integers/http"
	"github.com/aaronland/go-artisanal-integers/service"
	"github.com/aaronland/go-http-server"
	"log"
	gohttp "net/http"
	"net/url"
)

func NewServerApplicationFlags() *flag.FlagSet {

	fs := NewFlagSet("server")

	AssignCommonFlags(fs)

	fs.Int("set-last-int", 0, "Set the last known integer.")
	fs.Int("set-offset", 0, "Set the offset used to mint integers.")
	fs.Int("set-increment", 0, "Set the increment used to mint integers.")

	return fs
}

type ServerApplication struct {
	Application
}

func NewServerApplication(ctx context.Context, uri string) (Application, error) {

	a := &ServerApplication{}

	return a, nil
}

func (s *ServerApplication) Run(ctx context.Context, fl *flag.FlagSet) error {

	if !fl.Parsed() {
		ParseFlags(fl)
	}

	last, _ := IntVar(fl, "set-last-int")
	offset, _ := IntVar(fl, "set-last-offset")
	increment, _ := IntVar(fl, "set-last-increment")

	service_uri := "simple://"
	engine_uri := "memory://"

	server_uri := "http://localhost:8080"

	svc_uri, err := url.Parse(service_uri)

	if err != nil {
		return nil
	}

	svc_params := url.Values{}
	svc_params.Set("engine", engine_uri)

	svc_uri.RawQuery = svc_params.Encode()

	svc, err := service.NewService(ctx, svc_uri.String())

	if err != nil {
		return err
	}

	if last != 0 {

		err := s.engine.SetLastInt(int64(last))

		if err != nil {
			return err
		}
	}

	if increment != 0 {

		err := s.engine.SetIncrement(int64(increment))

		if err != nil {
			return err
		}
	}

	if offset != 0 {

		err := s.engine.SetOffset(int64(offset))

		if err != nil {
			return err
		}
	}

	//

	svr, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return err
	}

	mux := gohttp.NewServeMux()

	//

	integer_handler, err := http.IntegerHandler(svc)

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

	log.Println("Listen on", svr.Address())

	return svr.ListenAndServe(svc)
}
