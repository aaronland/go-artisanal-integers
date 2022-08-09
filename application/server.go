package application

import (
	"context"
	"flag"
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

		err := svc.SetLastInt(int64(last))

		if err != nil {
			return err
		}
	}

	if increment != 0 {

		err := svc.SetIncrement(int64(increment))

		if err != nil {
			return err
		}
	}

	if offset != 0 {

		err := svc.SetOffset(int64(offset))

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
		return err
	}

	integer_path := "/"

	mux.Handle(integer_path, integer_handler)

	log.Println("Listen on", svr.Address())

	return svr.ListenAndServe(ctx, mux)
}
