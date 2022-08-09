package service

import (
	"context"
	"errors"
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-artisanal-integers/engine"
	"net/url"
)

type SimpleService struct {
	Service
	engine engine.Engine
}

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "simple", NewSimpleService)

	if err != nil {
		panic(err)
	}
}

func NewSimpleService(ctx context.Context, uri string) (Service, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	engine_uri := q.Get("engine")

	if engine_uri == "" {
		return nil, errors.New("Missing ?engine= parameter")
	}

	eng, err := engine.NewEngine(ctx, engine_uri)

	if err != nil {
		return nil, err
	}

	svc := &SimpleService{
		engine: eng,
	}

	return svc, nil
}

func (svc *SimpleService) NextInt() (int64, error) {

	i, err := svc.engine.NextInt()

	if err != nil {
		return -1, err
	}

	if artisanalinteger.IsLondonInteger(i) {
		return svc.NextInt()
	}

	return i, nil
}

func (svc *SimpleService) LastInt() (int64, error) {
	return svc.engine.LastInt()
}

func (svc *SimpleService) SetLastInt(i int64) error {
	return svc.engine.SetLastInt(i)
}

func (svc *SimpleService) SetOffset(i int64) error {
	return svc.engine.SetOffset(i)
}

func (svc *SimpleService) SetIncrement(i int64) error {
	return svc.engine.SetIncrement(i)
}
