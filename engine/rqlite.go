package engine

// https://github.com/rqlite/rqlite/blob/master/doc/DATA_API.md

import (
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"net/http"
	"strconv"
	"sync"
)

type RqliteEngine struct {
	artisanalinteger.Engine
	endpoint  string
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
	LastInsertID int64
	RowsAffected int64
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

	sql := fmt.Sprintf("SELECT MAX(id) FROM %s", eng.key)

	results, err := eng.query(sql)

	if err != nil {
		return -1, err
	}

	r := results.Results[0]
	str_i := r.Values[0]
	i, err := strconv.ParseInt(str_i, 10, 64)

	if err != nil {
		return -1, err
	}

	return i, nil
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
		endpoint:  dsn,
		key:       "integers",
		increment: 2,
		offset:    1,
		mu:        mu,
		client:    client,
	}

	return &eng, nil
}
