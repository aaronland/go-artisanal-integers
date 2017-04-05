package engine

// https://github.com/rqlite/rqlite/blob/master/doc/DATA_API.md
// https://sqlite.org/autoinc.html

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thisisaaronland/go-artisanal-integers"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type RqliteEngine struct {
	artisanalinteger.Engine
	leader    string
	peers     []string
	key       string
	increment int64
	offset    int64
	mu        *sync.Mutex
	client    *http.Client
}

type Status struct {
	Build   BuildStatus   `json:"build"`
	HTTP    HTTPStatus    `json:"http"`
	Node    NodeStatus    `json:"node"`
	Runtime RuntimeStatus `json:"runtime"`
	Store   StoreStatus   `json:"store"`
}

type BuildStatus struct {
	Branch    string `json:"branch"`
	BuildTime string `json:"build_time"`
	Commit    string `json:"commit"`
	Version   string `json:"version"`
}

type HTTPStatus struct {
	Addr     string `json:"addr"`
	Auth     string `json:"auth"`
	Redirect string `json:"redirect"`
}

type NodeStatus struct {
	StartTime string `json:"start_time"`
	Uptime    string `json:"uptime"`
}

type RuntimeStatus struct {
	GOARCH       string `json:"GOARCH"`
	GOMAXPROCS   int    `json:"GOMAXPROCS"`
	GOOS         string `json:"GOOS"`
	NumCPU       int    `json:"numCPU"`
	NumGoRoutine int    `json:"numGoroutine"`
	Version      string `json:"version"`
}

// curl localhost:4003/status
type StoreStatus struct {
	Addr         string       `json:"addr"`
	ApplyTimeout string       `json:"apply_timeout"`
	DbConf       DbConfStatus `json:"db_conf"`
	Dir          string       `json:"dir"`
	Leader       string       `json:"leader"`
	Meta         MetaStatus   `json:"meta"`
	Peers        []string     `json:"peers"`
	Raft         RaftStatus   `json:"raft"`
	Sqlite3      SqliteStatus `json:"sqlite3"`
}

type DbConfStatus struct {
	DSN    string `json:"DSN"`
	Memory bool   `json:"Memory"`
}

type MetaStatus struct {
	APIPeers APIPeersStatus `json:"APIPeers"`
}

type APIPeersStatus map[string]string

func (p APIPeersStatus) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type RaftStatus struct {
	AppliedIndex      string `json:"applied_index"`
	CommitIndex       string `json:"commit_index"`
	FsmPending        string `json:"fsm_pending"`
	LastContact       string `json:"last_contact"`
	LastLogIndex      string `json:"last_log_index"`
	LastLogTerm       string `json:"last_log_term"`
	LastSnapshotIndex string `json:"last_snapshot_index"`
	LastSnapshotTerm  string `json:"last_snapshot_term"`
	NumPeers          string `json:"num_peers"`
	State             string `json:"state"`
	Term              string `json:"term"`
}

type SqliteStatus struct {
	DNS           string `json:"DNS"`
	FkConstraints string `json:"fk_constraints"`
	Path          string `json:"memory"`
	Version       string `json:"version"`
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

func get_rqlite_status(endpoint string) (*Status, error) {

	req, err := http.NewRequest("GET", endpoint+"/status", nil)

	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	rsp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	var status Status

	err = json.Unmarshal(body, &status)

	if err != nil {
		return nil, err
	}

	return &status, nil
}

func get_rqlite_peers(endpoint string) (string, []string, error) {

	var leader string
	var peers []string

	status, err := get_rqlite_status(endpoint)

	if err != nil {
		return leader, peers, err
	}

	store := status.Store
	meta := store.Meta

	/*

			See what's going on here? We want to point to the thing on port 4003
			and _not_ port 4004. It's weird. It appears to an Rqlite thing not a
			fast thing. Maybe? Dunno... (20170330/thisisaaronland)

		        "leader": "127.0.0.1:4004",
		        "meta": {
		            "APIPeers": {
		                "127.0.0.1:4002": "localhost:4001",
		                "127.0.0.1:4004": "localhost:4003"
		            }
		        },
		        "peers": [
		            "127.0.0.1:4004",
		            "127.0.0.1:4002"
		        ],

	*/

	leader, ok := meta.APIPeers[store.Leader]

	if !ok {

		if leader == "" {
			msg := fmt.Sprintf("Could not find entry for '%s' in API peers (%s)", store.Leader, meta.APIPeers)
			return leader, peers, errors.New(msg)
		}
	}

	leader = fmt.Sprintf("http://%s", leader)

	for _, host := range meta.APIPeers {

		p := fmt.Sprintf("http://%s", host)

		if p != endpoint {
			peers = append(peers, p)
		}
	}

	return leader, peers, nil
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

	if len(results.Results) == 0 {
		msg := fmt.Sprintf("Null results '%s' (%s)", results.Results, eng.leader)
		return -1, errors.New(msg)
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

	req, err := http.NewRequest("GET", eng.leader+"/db/query", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	params := url.Values{}
	params.Set("q", sql)

	req.URL.RawQuery = (params).Encode()

	// So this req.URL will fail with a 400 error when you use a host that
	// is **not** currently the leader... rather than something that will
	// return a 301. Even though the exact same URL sent via curl works fine.
	// Computers amirite? (20170330/thisisaaronland)

	rsp, err := eng.do(req)

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

	req, err := http.NewRequest("POST", eng.leader+"/db/execute", buf)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	rsp, err := eng.do(req)

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

func (eng *RqliteEngine) do(req *http.Request) (*http.Response, error) {

	// Hack - see below

	var b bytes.Buffer
	var buf *bytes.Buffer

	if req.Body != nil {

		wr := bufio.NewWriter(&b)

		io.Copy(wr, req.Body)

		buf = bytes.NewBuffer(b.Bytes())
		req.Body = ioutil.NopCloser(buf)
	}

	var retry bool
	var retry_leader string

	rsp, err := eng.client.Do(req)

	if err != nil && len(eng.peers) > 0 {

		var new_leader string
		var new_peers []string

		keep_trying := true

		counter := 0
		max := len(eng.peers) * 100

		eng.mu.Lock()

		for {

			counter += 1

			if counter >= max {
				fmt.Fprintf(os.Stderr, "Couldn't find new leader after %d tries so giving up\n", max)
				break
			}

			fmt.Fprintf(os.Stderr, "Couldn't connect to leader so trying to see if the peers are rebalancing themselves (%d/%d)...\n", counter, max)

			for _, pr := range eng.peers {

				leader, peers, err := get_rqlite_peers(pr)

				if err != nil {
					keep_trying = false
					break
				}

				if leader != eng.leader {

					new_leader = leader
					new_peers = peers

					keep_trying = false
					retry = true
				}
			}

			if !keep_trying {
				break
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}

		if retry {

			eng.leader = new_leader
			eng.peers = new_peers

			retry_leader = eng.leader
			retry = true
		}

		eng.mu.Unlock()

	} else if rsp.StatusCode == 301 {

		rsp.Body.Close()

		location := rsp.Header.Get("Location")
		leader, err := url.Parse(location)

		if err != nil {
			return nil, err
		}

		new_leader := fmt.Sprintf("%s://%s", leader.Scheme, leader.Host)
		eng.leader = new_leader

		retry = true
		retry_leader = new_leader

	} else {
		retry = false
	}

	if retry {

		new_url, err := url.Parse(retry_leader + req.URL.Path)

		if err != nil {
			return nil, err
		}

		req.URL = new_url

		// Hack - see below

		if req.Body != nil {
			buf = bytes.NewBuffer(b.Bytes())
			req.Body = ioutil.NopCloser(buf)
		}

		// FIX ME: why is req.Body being closed even though it's a *bytes.Buffer?
		// https://golang.org/pkg/net/http/#NewRequest
		// Answer: Because an older version of Go

		return eng.do(req)
	}

	return rsp, err
}

func NewRqliteEngine(dsn string) (*RqliteEngine, error) {

	leader, peers, err := get_rqlite_peers(dsn)

	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	mu := new(sync.Mutex)

	eng := RqliteEngine{
		leader:    leader,
		peers:     peers,
		key:       "integers",
		increment: 2,
		offset:    1,
		mu:        mu,
		client:    client,
	}

	go func() {

		timer := time.NewTimer(time.Second * 1).C
		done := make(chan bool)

		for {
			select {
			case <-timer:

				leader, peers, err := get_rqlite_peers(eng.leader)

				if err != nil {
					done <- true
				}

				if leader != eng.leader {
					eng.mu.Lock()
					eng.leader = leader
					eng.peers = peers
					eng.mu.Unlock()
				}

			case <-done:
				break
			default:
				//
			}
		}
	}()

	return &eng, nil
}
