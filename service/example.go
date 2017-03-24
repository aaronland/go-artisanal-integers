package service

import (
	"github.com/thisisaaronland/go-artisanal-integers"
)

type ExampleService struct {
	artisanalinteger.Service
	engine 	*artisanalinteger.Engine
}

func NewExampleService (eng *artisanalinteger.Engine) (*ExampleService, error) {

     srv := ExampleService{
     	engine: eng,
     }

     return &src, nil
}

func (srv *ExampleService) NextId() (int64, error) {
     return srv.engine.NextId()
}

func (srv *ExampleService) LastId() (int64, error) {
     return srv.engine.LastId()
}