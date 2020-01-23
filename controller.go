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

type Method interface {
}
type PostMethod Method
