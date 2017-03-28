package engine

import (
	"errors"
	"github.com/thisisaaronland/go-artisanal-integers"
	"sync"
)

type RqliteEngine struct {
	artisanalinteger.Engine
	key       string
	increment int64
	offset    int64
	mu        *sync.Mutex
}

func (eng *RqliteEngine) SetLastInt(i int64) error {
	return errors.New("Please implement me")
}

func (eng *RqliteEngine) SetKey(k string) error {
	return errors.New("Please implement me")
}

func (eng *RqliteEngine) SetOffset(i int64) error {
	return errors.New("Please implement me")
}

func (eng *RqliteEngine) SetIncrement(i int64) error {
	return errors.New("Please implement me")
}

func (eng *RqliteEngine) NextInt() (int64, error) {
	return -1, errors.New("Please implement me")
}

func (eng *RqliteEngine) LastInt() (int64, error) {
	return -1, errors.New("Please implement me")
}

func NewRqliteEngine(dsn string) (*RqliteEngine, error) {

	mu := new(sync.Mutex)

	eng := RqliteEngine{
		key:       "integers",
		increment: 2,
		offset:    1,
		mu:        mu,
	}

	return &eng, nil
}
