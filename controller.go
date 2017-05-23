package gerry

import (
	"github.com/gorilla/mux"
	"net/http"
)

type HttpAction func(ctx *Context)

type Controller interface {
	InitActions(router *mux.Router)
}

func NewAction(action HttpAction) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		context := NewContext(rw, req)
		action(context)
	}
}
