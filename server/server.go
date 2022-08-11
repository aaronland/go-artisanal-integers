package server

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type Server interface {
	Address() string
	ListenAndServe(context.Context, ...interface{}) error
}

type ServerInitializeFunc func(ctx context.Context, uri string) (Server, error)

var servers roster.Roster

func ensureServerRoster() error {

	if servers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		servers = r
	}

	return nil
}

func RegisterServer(ctx context.Context, scheme string, f ServerInitializeFunc) error {

	err := ensureServerRoster()

	if err != nil {
		return err
	}

	return servers.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureServerRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range servers.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := servers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(ServerInitializeFunc)
	s, err := f(ctx, uri)

	if err != nil {
		return nil, err
	}

	return s, nil
}
