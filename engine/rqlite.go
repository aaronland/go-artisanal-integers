package engine

// https://github.com/rqlite/rqlite/blob/master/doc/DATA_API.md

import (
	"errors"
	"github.com/thisisaaronland/go-artisanal-integers"
	"net/http"
	"sync"
)

type RqliteEngine struct {
	artisanalinteger.Engine
	key       string
	increment int64
	offset    int64
	mu        *sync.Mutex
	client    *http.Client
}

type QueryTime float64

type QueryResults struct {
	Results []QueryResult
	Time    QueryTime
}

type QueryResult struct {
	Columns []string
	Types   []string
	Values  []string
	Time    QueryTime
}

type ExecuteResults struct {
	Results []ExecuteResult
	Time    QueryTime
}

type ExecuteResult struct {
	LastInsertID int
	RowsAffected int
	Time         QueryTime
}

func (eng *RqliteEngine) SetLastInt(i int64) error {
	return errors.New("Please implement me")
}

func (eng *RqliteEngine) SetKey(k string) error {
	eng.key = k
	return nil
}

func (eng *RqliteEngine) SetOffset(i int64) error {
	eng.offset = i
	return nil
}

func (eng *RqliteEngine) SetIncrement(i int64) error {
	eng.increment = i
	return nil
}

func (eng *RqliteEngine) NextInt() (int64, error) {
	return -1, errors.New("Please implement me")
}

func (eng *RqliteEngine) LastInt() (int64, error) {
	return -1, errors.New("Please implement me")
}

func (eng *RqliteEngine) query(sql string) (*QueryResults, error) {
	return nil, errors.New("Please implement me")
}

func (eng *RqliteEngine) execute(sql string) (*ExecuteResults, error) {
	return nil, errors.New("Please implement me")
}

func NewRqliteEngine(dsn string) (*RqliteEngine, error) {

	client := new(http.Client)
	mu := new(sync.Mutex)

	eng := RqliteEngine{
		key:       "integers",
		increment: 2,
		offset:    1,
		mu:        mu,
		client:    client,
	}

	return &eng, nil
}
