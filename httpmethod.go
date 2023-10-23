package cango

import (
	"net/http"
)

type GetMethod interface {
	cangoHttpGet()
}

type PostMethod interface {
	cangoHttpPost()
}

type PutMethod interface {
	cangoHttpPut()
}

type DeleteMethod interface {
	cangoHttpDelete()
}

type PatchMethod interface {
	cangoHttpPatch()
}

type OptionsMethod interface {
	cangoHttpOptions()
}

type HeadMethod interface {
	cangoHttpHead()
}

type TraceMethod interface {
	cangoHttpTrace()
}

func httpMethods(ps any) (methods []string) {
	if _, ok := ps.(GetMethod); ok {
		methods = append(methods, http.MethodGet)
	}
	if _, ok := ps.(PostMethod); ok {
		methods = append(methods, http.MethodPost)
	}
	if _, ok := ps.(PutMethod); ok {
		methods = append(methods, http.MethodPut)
	}
	if _, ok := ps.(DeleteMethod); ok {
		methods = append(methods, http.MethodDelete)
	}
	if _, ok := ps.(PatchMethod); ok {
		methods = append(methods, http.MethodPatch)
	}
	if _, ok := ps.(OptionsMethod); ok {
		methods = append(methods, http.MethodOptions)
	}
	if _, ok := ps.(HeadMethod); ok {
		methods = append(methods, http.MethodHead)
	}
	if _, ok := ps.(TraceMethod); ok {
		methods = append(methods, http.MethodTrace)
	}
	// 默认为Get
	if len(methods) == 0 {
		methods = append(methods, http.MethodGet)
	}
	return
}
