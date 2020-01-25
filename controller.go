package cango

import (
	"net/http"
)

type URI interface {
	Request() *http.Request
}

type context struct {
	request *http.Request
}

func (c *context) Request() *http.Request {
	return c.request
}

type ModelView struct {
	Tpl   string
	Model interface{}
}

type Redirect struct {
	Url  string
	Code int
}

type httpMethod interface {
}

/*
	methods := []string{"Get", "Post", "Head", "Put", "Patch", "Delete", "Options", "Trace"}
	for _, m := range methods {
		fmt.Printf("type %sMethod httpMethod\n", m)
	}
*/

type GetMethod httpMethod
type HeadMethod httpMethod
type PostMethod httpMethod
type PutMethod httpMethod
type PatchMethod httpMethod
type DeleteMethod httpMethod
type OptionsMethod httpMethod
type TraceMethod httpMethod
