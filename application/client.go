package application

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/client"
	"log"
	"net/url"
)

func NewClientApplicationFlags() *flag.FlagSet {

	fs := NewFlagSet("client")

	AssignCommonFlags(fs)

	return fs
}

type ClientApplication struct {
	Application
}

func NewClientApplication(ctx context.Context, uri string) (Application, error) {

	c := &ClientApplication{}
	return c, nil
}

func (c *ClientApplication) Run(ctx context.Context, fl *flag.FlagSet) error {

	if !fl.Parsed() {
		ParseFlags(fl)
	}

	client_uri := "http://localhost:8080"

	cl, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return err
	}

	i, err := cl.NextInt()

	if err != nil {
		return err
	}

	log.Println(i)
	return nil
}
