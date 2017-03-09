package engine

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/thisisaaronland/go-artisanal-integers"
	"strconv"
	"sync"
)

type RedisEngine struct {
	artisanalinteger.Engine
	pool   *redis.Pool
	key    string
	incrby int
	mu     *sync.Mutex
}

func (eng *SummitDBEngine) Set(i int64) error {

	eng.mu.Lock()
	defer m.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", m.key, i)
	return err
}

func (eng *RedisEngine) Max() (int64, error) {

	eng.mu.Lock()
	defer m.mu.Unlock()

	conn := eng.pool.Get()
	defer conn.Close()

	redis_rsp, err := conn.Do("GET", m.key)

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

func (eng *RedisEngine) Next() (int64, error) {

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

func NewRedisEngine(host string, port int, key string, incrby int) (*SummitDBEngine, error) {

	db_endpoint := fmt.Sprintf("%s:%d", host, port)

	pool := &redis.Pool{
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial("tcp", db_endpoint)

			if err != nil {
				return nil, err
			}

			return c, err
		},
	}

	mu := new(sync.Mutex)

	eng := RedisEngine{
		pool:   pool,
		key:    key,
		incrby: incrby,
		mu:     mu,
	}

	return &eng, nil
}
