package engine

// https://github.com/rqlite/rqlite/blob/master/doc/DATA_API.md

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"io/ioutil"
	"net/http"
	"net/url"
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

	params := url.Values{}
	params.Set("q", sql)

	req, err := http.NewRequest("GET", eng.endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = (params).Encode()

	rsp, err := eng.client.Do(req)

	if err != nil {
		msg := fmt.Sprintf("HTTP request failed: %s", err.Error())
		return nil, errors.New(msg)
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	var results QueryResults

	err = json.Unmarshal(body, &results)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (eng *RqliteEngine) execute(sql string) (*ExecuteResults, error) {

	// FIX ME - POST body data...

	params := url.Values{}
	params.Set("q", sql)

	req, err := http.NewRequest("POST", eng.endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = (params).Encode()

	rsp, err := eng.client.Do(req)

	if err != nil {
		msg := fmt.Sprintf("HTTP request failed: %s", err.Error())
		return nil, errors.New(msg)
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	var results ExecuteResults

	err = json.Unmarshal(body, &results)

	if err != nil {
		return nil, err
	}

	return &results, nil
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
