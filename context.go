package gerry

import (
	"context"
	"github.com/unrolled/render"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewContext(resp http.ResponseWriter, req *http.Request) *Context {
	r := render.New(render.Options{
		Directory:  "templates",
		IndentJSON: true,
	})

	ctx := NewContextWithRender(resp, req, r)
	ctx.Request.ParseForm()
	ctx.params = ctx.Request.Form

	return ctx
}

func NewContextWithRender(resp http.ResponseWriter, req *http.Request, rndr *render.Render) *Context {
	return &Context{req, resp, rndr, nil}
}

// type ContextHandlerFunc func(ctx *Context, next http.HandlerFunc)

type ContextHandler func(ctx *Context, next http.HandlerFunc)

func (h ContextHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := NewContext(rw, r)
	h(ctx, next)
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Render   *render.Render
	params   url.Values
}

func (c *Context) Merge(vals ...url.Values) url.Values {
	merge := url.Values{}

	for _, val := range vals {
		for key, value := range val {
			merge[key] = value
		}
	}

	return merge
}

func (c *Context) WithContext(context context.Context) {
	c.Request = c.Request.WithContext(context)
}

func (c *Context) WithContextValue(key, value interface{}) {
	context := context.WithValue(c.Request.Context(), key, value)
	c.WithContext(context)
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}

func (c *Context) ContextValue(key interface{}) interface{} {
	return c.Request.Context().Value(key)
}

func (c *Context) Paths() []string {
	return strings.Split(c.Request.URL.Path, "/")
}

func (c *Context) Path() string {
	return c.Request.URL.Path
}

func (c *Context) URL() *url.URL {
	return c.Request.URL
}

func (c *Context) Exec(f http.HandlerFunc) {
	f(c.Response, c.Request)
}

func (c *Context) Redirect(urlStr string, code int) {
	http.Redirect(c.Response, c.Request, urlStr, code)
}

func (c *Context) ParamInt(key string) int {
	i, _ := strconv.Atoi(c.Param(key))

	return i
}

func (c *Context) ParamInt64(key string) int64 {
	i, _ := strconv.ParseInt(c.Param(key), 10, 64)

	return i
}

func (c *Context) ParamFloat(key string) float64 {
	f, _ := strconv.ParseFloat(c.Param(key), 64)

	return f
}

func (c *Context) Param(key string) string {
	val, _ := c.Get(key)

	return val
}

func (c *Context) Get(key string) (string, bool) {
	if values, ok := c.params[key]; ok {
		return values[0], ok
	}

	if val := c.Request.FormValue(key); val != "" {
		return val, true
	}

	if val := c.Request.PostFormValue(key); val != "" {
		return val, true
	}

	return "", false
}

func (c *Context) GetStrings(key string) ([]string, bool) {
	values, ok := c.params[key]

	return values, ok
}

func (c *Context) Text(v string) {
	c.TextWithCode(http.StatusOK, v)
}

func (c *Context) JSON(v interface{}) {
	c.JSONWithCode(http.StatusOK, v)
}

func (c *Context) HTML(name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	c.HTMLWithCode(http.StatusOK, name, binding, htmlOpt...)
}

func (c *Context) Data(v []byte) {
	c.DataWithCode(http.StatusOK, v)
}

func (c *Context) Error(v string) {
	c.ErrorWithCode(http.StatusInternalServerError, v)
}

func (c *Context) jsonText(code int, v string) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.TextWithCode(code, v)
}

func (c *Context) TextWithCode(code int, v string) {
	c.Render.Text(c.Response, code, v)
}

func (c *Context) JSONWithCode(code int, v interface{}) {
	switch v.(type) {
	case string:
		c.jsonText(code, v.(string))
	default:
		c.Render.JSON(c.Response, code, v)
	}
}

func (c *Context) HTMLWithCode(code int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	c.Render.HTML(c.Response, code, name, binding, htmlOpt...)
}

func (c *Context) DataWithCode(code int, v []byte) {
	c.Render.Data(c.Response, code, v)
}

func (c *Context) ErrorWithCode(code int, v string) {
	http.Error(c.Response, v, code)
}
