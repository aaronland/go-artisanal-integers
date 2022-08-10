package service

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Service interface {
	NextInt(context.Context) (int64, error)
	LastInt(context.Context) (int64, error)
	SetLastInt(context.Context, int64) error
	SetOffset(context.Context, int64) error
	SetIncrement(context.Context, int64) error
	Close(context.Context) error
}

type ServiceInitializeFunc func(ctx context.Context, uri string) (Service, error)

var services roster.Roster

func ensureServiceRoster() error {

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

	err := ensureServiceRoster()

	if err != nil {
		return err
	}

	return services.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureServiceRoster()

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
	s, err := f(ctx, uri)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func SetParametersFromURI(ctx context.Context, s Service, uri string) error {

	u, err := url.Parse(uri)

	if err != nil {
		return fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	str_offset := q.Get("offset")
	str_increment := q.Get("increment")
	str_last := q.Get("last-int")

	if str_offset != "" {

		offset, err := strconv.ParseInt(str_offset, 10, 64)

		if err != nil {
			return fmt.Errorf("Invalid ?offset= parameter, %w", err)
		}

		err = s.SetOffset(ctx, offset)

		if err != nil {
			return fmt.Errorf("Failed to set offset, %w", err)
		}
	}

	if str_increment != "" {

		increment, err := strconv.ParseInt(str_increment, 10, 64)

		if err != nil {
			return fmt.Errorf("Invalid ?increment= parameter, %w", err)
		}

		err = s.SetIncrement(ctx, increment)

		if err != nil {
			return fmt.Errorf("Failed to set increment, %w", err)
		}
	}

	if str_last != "" {

		last, err := strconv.ParseInt(str_last, 10, 64)

		if err != nil {
			return fmt.Errorf("Invalid ?last= parameter, %w", err)
		}

		err = s.SetLastInt(ctx, last)

		if err != nil {
			return fmt.Errorf("Failed to set last, %w", err)
		}
	}

	return nil
}
