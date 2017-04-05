package util

import (
	"errors"
	"github.com/thisisaaronland/go-artisanal-integers"
	"github.com/thisisaaronland/go-artisanal-integers/engine"
)

func NewArtisanalEngine(db string, dsn string) (artisanalinteger.Engine, error) {

	var eng artisanalinteger.Engine
	var err error

	switch db {

	case "memory":

		eng, err = engine.NewMemoryEngine(dsn)

	case "mysql":

		eng, err = engine.NewMySQLEngine(dsn)

	case "redis":

		if dsn == "" {
			dsn = "localhost:6379"
		}

		eng, err = engine.NewRedisEngine(dsn)

	case "rqlite":

		if dsn == "" {
			dsn = "http://localhost:4001"
		}

		eng, err = engine.NewRqliteEngine(dsn)

	case "summitdb":

		if dsn == "" {
			dsn = "localhost:7481"
		}

		eng, err = engine.NewSummitDBEngine(dsn)

	default:
		return nil, errors.New("Invalid engine")
	}

	if err != nil {
		return nil, err
	}

	return eng, nil
}
