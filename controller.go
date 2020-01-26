package cango

import (
	"net/http"
)

type WebRequest struct {
	*http.Request
	localHolder map[string]interface{}
}

type URI interface {
	Request() WebRequest
}

type context struct {
	WebRequest
}

func (c *context) Request() WebRequest {
	return c.WebRequest
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
