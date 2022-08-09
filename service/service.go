package service

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type Service interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	//SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
}

type ServiceInitializeFunc func(ctx context.Context, uri string) (Service, error)

var services roster.Roster

func ensureSpatialRoster() error {

	if services == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		services = r
	}

	return nil
}

func RegisterService(ctx context.Context, scheme string, f ServiceInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return services.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range services.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewService(ctx context.Context, uri string) (Service, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := services.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(ServiceInitializeFunc)
	return f(ctx, uri)
}
