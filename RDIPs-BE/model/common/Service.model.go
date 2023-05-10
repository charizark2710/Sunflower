package model

import (
	"context"
	"net/http"
	"sync"
)

type ServiceModel struct {
	param  map[string]string
	query  map[string]string
	Body   []byte
	Header http.Header
}

type ServiceContext struct {
	Ctx context.Context
	Mu  sync.Mutex
	ServiceModel
}

func (sctx *ServiceContext) InitParamsAndQueries() {
	sctx.param = make(map[string]string)
	sctx.query = make(map[string]string)
}

func (ctx *ServiceContext) SetParam(key string, value string) {
	ctx.param[key] = value
}

func (ctx *ServiceContext) SetQuery(key string, value string) {
	ctx.query[key] = value
}

func (ctx *ServiceContext) Param(key string) string {
	return ctx.param[key]
}

func (ctx *ServiceContext) Query(key string) string {
	return ctx.query[key]
}
