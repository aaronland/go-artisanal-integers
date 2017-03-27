package engine

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/thisisaaronland/go-artisanal-integers"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func make_pool(dsn string) (*redis.Pool, error) {

	pool := &redis.Pool{
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {

			c, err := redis.DialURL(dsn)

			if err != nil {
				return nil, err
			}

			return c, err
		},
	}

	return pool, nil
}

type SummitDBEngine struct {
	artisanalinteger.Engine
	pool      *redis.Pool
	peers     []string
	key       string
	increment int64
	offset    int64
	mu        *sync.Mutex
}

func (eng *SummitDBEngine) SetLastInt(i int64) error {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", eng.key, i)
	return err
}

func (eng *SummitDBEngine) SetKey(k string) error {
	eng.key = k
	return nil
}

func (eng *SummitDBEngine) SetOffset(i int64) error {
	eng.offset = i
	return nil
}

func (eng *SummitDBEngine) SetIncrement(i int64) error {
	eng.increment = i
	return nil
}

func (eng *SummitDBEngine) LastInt() (int64, error) {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	redis_rsp, err := conn.Do("GET", eng.key)

	if err != nil {
		return -1, err
	}

	b, err := redis.Bytes(redis_rsp, nil)

	if err != nil {
		return -1, err
	}

	i, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func (eng *SummitDBEngine) NextInt() (int64, error) {

	i, err := eng.nextInt()

	if err != nil {

		if strings.HasPrefix(err.Error(), "TRY") {

			parsed := strings.Split(err.Error(), " ")
			dsn := fmt.Sprintf("redis://%s", parsed[1])

			fmt.Fprintf(os.Stderr, "summitdb told me to try %s instead, so here we go...\n", dsn)

			eng.mu.Lock()

			// See the way we're explicitly unlocking the mutex rather
			// than defer-ing it on exit? Yes, that because we are potentially
			// going to call ourselves recursively here which does not invoke
			// the defer robot (20170327/thisisaaronland)

			pool, err := make_pool(dsn)

			if err != nil {
				eng.mu.Unlock()
				return -1, err
			}

			err = eng.pool.Close()

			if err != nil {
				eng.mu.Unlock()
				return -1, err
			}

			eng.pool = pool
			eng.mu.Unlock()

			return eng.NextInt()
		}

		return -1, err
	}

	return i, nil
}

func (eng *SummitDBEngine) nextInt() (int64, error) {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	redis_rsp, err := conn.Do("INCRBY", eng.key, eng.increment)

	if err != nil {
		return -1, err
	}

	i, err := redis.Int64(redis_rsp, nil)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func (eng *SummitDBEngine) get_peers() error {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	redis_rsp, err := conn.Do("RAFTPEERS")

	if err != nil {
		return err
	}

	log.Println(redis_rsp)

	return nil
}

func NewSummitDBEngine(dsn string) (*SummitDBEngine, error) {

	pool, err := make_pool(dsn)

	if err != nil {
		return nil, err
	}

	peers := make([]string, 0)

	mu := new(sync.Mutex)

	eng := SummitDBEngine{
		pool:      pool,
		peers:     peers,
		key:       "integers",
		increment: 2,
		offset:    1,
		mu:        mu,
	}

	eng.get_peers()

	return &eng, nil
}
