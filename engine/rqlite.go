package engine

// https://github.com/rqlite/rqlite/blob/master/doc/DATA_API.md
// https://sqlite.org/autoinc.html

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"io/ioutil"
	_ "log"
	"net/http"
	"net/url"
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
	Results []QueryResult `json:"results"`
	Time    QueryTime     `json:"time"`
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Types   []string        `json:"types"`
	Values  [][]interface{} `json:"values"` // really we only care about int64s but just in case...
	Time    QueryTime       `json:"time"`
	Error   string          `json:"error"`
}

func (r *QueryResults) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

type ExecuteResults struct {
	Results []ExecuteResult `json:"results"`
	Time    QueryTime       `json:"time"`
}

func (r *ExecuteResults) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

type ExecuteResult struct {
	LastInsertId int64     `json:"last_insert_id"`
	RowsAffected int64     `json:"rows_affected"`
	Time         QueryTime `json:"time"`
	Error        string    `json:"error"`
}

func (eng *RqliteEngine) SetLastInt(i int64) error {

	last, err := eng.LastInt()

	if err != nil {
		return err
	}

	if i < last {
		return errors.New("integer value too small")
	}

	// https://stackoverflow.com/questions/692856/set-start-value-for-autoincrement-in-sqlite

	sql := fmt.Sprintf("UPDATE sqlite_sequence SET seq=%d WHERE name='%s'", i, eng.key)

	_, err = eng.execute(sql)
	return err
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

	sql := fmt.Sprintf("REPLACE INTO %s (stub) VALUES ('a')", eng.key)

	results, err := eng.execute(sql)

	if err != nil {
		return -1, err
	}

	r := results.Results[0]

	if r.Error != "" {
		return -1, errors.New(r.Error)
	}

	i := r.LastInsertId
	return i, nil
}

func (eng *RqliteEngine) LastInt() (int64, error) {

	sql := fmt.Sprintf("SELECT MAX(id) FROM %s", eng.key)

	results, err := eng.query(sql)

	if err != nil {
		return -1, err
	}

	r := results.Results[0]

	if r.Error != "" {
		return -1, errors.New(r.Error)
	}

	values := r.Values[0]

	i := values[0].(float64)
	return int64(i), nil
}

func (eng *RqliteEngine) query(sql string) (*QueryResults, error) {

	params := url.Values{}
	params.Set("q", sql)

	req, err := http.NewRequest("GET", eng.endpoint+"/db/query", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
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

	q := []string{sql}

	b, err := json.Marshal(q)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", eng.endpoint+"/db/execute", buf)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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
