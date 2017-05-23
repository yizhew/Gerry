package gerry

import (
	log "github.com/Sirupsen/logrus"
	"github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
	"net/http"
	"time"
)

type Gerry struct {
	*negroni.Negroni
}

func New() *Gerry {
	return &Gerry{negroni.New()}
}

func (g *Gerry) UseContextHandler(handler *ContextHandler) {
	g.Use(handler)
}

func (g *Gerry) UseContextHandlerFunc(handlerFunc func(ctx *Context, next http.HandlerFunc)) {
	g.Use(ContextHandler(handlerFunc))
}

func (g *Gerry) ApplyRecovery(ehf func(err interface{})) {
	recovery := negroni.NewRecovery()

	if ehf == nil {
		g.Use(recovery)
		return
	}

	recovery.PrintStack = false
	recovery.ErrorHandlerFunc = ehf

	if ErrorLogger == nil {
		SetupErrorLogger("app.error")
	}
	recovery.Logger = ErrorLogger
	g.Use(recovery)
}

func (g *Gerry) ApplyLogger(l *log.Logger, beforeLogging func(entry *log.Entry, req *http.Request, remoteAddr string) *log.Entry, afterLogging func(entry *log.Entry, res negroni.ResponseWriter, latency time.Duration, name string) *log.Entry) {
	if l == nil {
		SetupLogger("gerry.log")
		l = Logger
	}

	nl := negronilogrus.NewMiddlewareFromLogger(l, "gerry")
	if beforeLogging != nil {
		// override the default Before
		nl.Before = beforeLogging
	}

	if afterLogging != nil {
		// override the default After
		nl.After = afterLogging
	}

	g.Use(nl)
}

func (g *Gerry) ApplyStatic(path string) {
	g.Use(negroni.NewStatic(http.Dir(path)))
}

func (g *Gerry) ApplyNoCache() {
	g.Use(newNoCacheMiddleware())
}
