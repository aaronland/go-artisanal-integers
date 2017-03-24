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

func (svc *ExampleService) NextId() (int64, error) {
	return svc.engine.NextId()
}

func (svc *ExampleService) LastId() (int64, error) {
	return svc.engine.LastId()
}
