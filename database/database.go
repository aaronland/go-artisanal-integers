package database

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Database interface {
	NextInt(context.Context) (int64, error)
	LastInt(context.Context) (int64, error)
	SetLastInt(context.Context, int64) error
	SetOffset(context.Context, int64) error
	SetIncrement(context.Context, int64) error
	Close(context.Context) error
}

type DatabaseInitializeFunc func(ctx context.Context, uri string) (Database, error)

var databases roster.Roster

func ensureSpatialRoster() error {

	if databases == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		databases = r
	}

	return nil
}

func RegisterDatabase(ctx context.Context, scheme string, f DatabaseInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return databases.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range databases.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewDatabase(ctx context.Context, uri string) (Database, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := databases.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}
	
	f := i.(DatabaseInitializeFunc)
	db, err := f(ctx, uri)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetParametersFromURI(ctx context.Context, db Database, uri string) error {

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

		err = db.SetOffset(ctx, offset)

		if err != nil {
			return fmt.Errorf("Failed to set offset, %w", err)
		}
	}

	if str_increment != "" {

		increment, err := strconv.ParseInt(str_increment, 10, 64)

		if err != nil {
			return fmt.Errorf("Invalid ?increment= parameter, %w", err)
		}

		err = db.SetIncrement(ctx, increment)

		if err != nil {
			return fmt.Errorf("Failed to set increment, %w", err)
		}
	}

	if str_last != "" {

		last, err := strconv.ParseInt(str_last, 10, 64)

		if err != nil {
			return fmt.Errorf("Invalid ?last= parameter, %w", err)
		}

		err = db.SetLastInt(ctx, last)

		if err != nil {
			return fmt.Errorf("Failed to set last, %w", err)
		}
	}
	

	return nil
}
