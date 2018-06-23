package service

import (
	"github.com/aaronland/go-artisanal-integers"
)

type SimpleService struct {
	artisanalinteger.Service
	engine artisanalinteger.Engine
}

func NewSimpleService(eng artisanalinteger.Engine) (*SimpleService, error) {

	svc := SimpleService{
		engine: eng,
	}

	return &svc, nil
}

func (svc *SimpleService) NextInt() (int64, error) {
	return svc.engine.NextInt()
}

func (svc *SimpleService) LastInt() (int64, error) {
	return svc.engine.LastInt()
}
