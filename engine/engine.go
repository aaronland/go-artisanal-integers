package engine

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type Engine interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
	Close() error
}

type EngineInitializeFunc func(ctx context.Context, uri string) (Engine, error)

var engines roster.Roster

func ensureSpatialRoster() error {

	if engines == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		engines = r
	}

	return nil
}

func RegisterEngine(ctx context.Context, scheme string, f EngineInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return engines.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range engines.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewEngine(ctx context.Context, uri string) (Engine, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := engines.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(EngineInitializeFunc)
	return f(ctx, uri)
}
