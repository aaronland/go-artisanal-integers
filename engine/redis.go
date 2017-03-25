package engine

import (
	"github.com/garyburd/redigo/redis"
	"github.com/thisisaaronland/go-artisanal-integers"
	"strconv"
	"sync"
)

type RedisEngine struct {
	artisanalinteger.Engine
	pool      *redis.Pool
	key       string
	offset    int64
	increment int64
	mu        *sync.Mutex
}

func (eng *RedisEngine) SetLastId(i int64) error {

	eng.mu.Lock()
	defer eng.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", eng.key, i)
	return err
}

func (eng *RedisEngine) SetKey(k string) error {
	eng.key = k
	return nil
}

func (eng *RedisEngine) SetOffset(i int64) error {
	eng.offset = i
	return nil
}

func (eng *RedisEngine) SetIncrement(i int64) error {
	eng.increment = i
	return nil
}

func (eng *RedisEngine) LastId() (int64, error) {

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

func (eng *RedisEngine) NextId() (int64, error) {

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

func NewRedisEngine(dsn string) (*RedisEngine, error) {

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

	eng := RedisEngine{
		pool:      pool,
		key:       "integers",
		increment: 2,
		offset:    1,
		mu:        mu,
	}

	return &eng, nil
}
