package engine

import (
	"github.com/garyburd/redigo/redis"
	"github.com/thisisaaronland/go-artisanal-integers"
	"strconv"
	"sync"
)

type SummitDBEngine struct {
	artisanalinteger.Engine
	pool   *redis.Pool
	key    string
	incrby int
	mu     *sync.Mutex
}

func (eng *SummitDBEngine) SetLastId(i int64) error {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", eng.key, i)
	return err
}

func (eng *SummitDBEngine) SetOffset(int64) error {
	return nil
}

func (eng *SummitDBEngine) SetIncrement(int64) error {
	return nil
}

func (eng *SummitDBEngine) LastId() (int64, error) {

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

func (eng *SummitDBEngine) NextId() (int64, error) {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	redis_rsp, err := conn.Do("INCRBY", eng.key, eng.incrby)

	if err != nil {
		return -1, err
	}

	i, err := redis.Int64(redis_rsp, nil)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func NewSummitDBEngine(dsn string) (*SummitDBEngine, error) {

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

	mu := new(sync.Mutex)

	eng := SummitDBEngine{
		pool:   pool,
		key:    "integers",
		incrby: 2,
		mu:     mu,
	}

	return &eng, nil
}
