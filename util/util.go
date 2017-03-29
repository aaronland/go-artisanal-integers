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

	case "redis":
		eng, err = engine.NewRedisEngine(dsn)
	case "rqlite":
		eng, err = engine.NewRqliteEngine(dsn)
	case "summitdb":
		eng, err = engine.NewSummitDBEngine(dsn)
	case "mysql":
		eng, err = engine.NewMySQLEngine(dsn)
	default:
		return nil, errors.New("Invalid engine")
	}

	if err != nil {
		return nil, err
	}

	return eng, nil
}
