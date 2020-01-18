package cango

import (
	"net/http"
)

const (
	GET     = http.MethodGet
	HEAD    = http.MethodHead
	POST    = http.MethodPost
	PUT     = http.MethodPut
	PATCH   = http.MethodPatch
	DELETE  = http.MethodDelete
	OPTIONS = http.MethodOptions
	TRACE   = http.MethodTrace
)
