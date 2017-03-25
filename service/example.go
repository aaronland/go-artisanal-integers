package service

import (
	"github.com/thisisaaronland/go-artisanal-integers"
)

type ExampleService struct {
	artisanalinteger.Service
	engine artisanalinteger.Engine
}

func NewExampleService(eng artisanalinteger.Engine) (*ExampleService, error) {

	svc := ExampleService{
		engine: eng,
	}

	return &svc, nil
}

func (svc *ExampleService) NextInt() (int64, error) {
	return svc.engine.NextInt()
}

func (svc *ExampleService) LastInt() (int64, error) {
	return svc.engine.LastInt()
}
