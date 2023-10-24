package cango

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type URI interface {
	ginContext() *gin.Context
	Context() *gin.Context
}

type defaultURI struct {
	ctx *gin.Context
}

func (uri *defaultURI) ginContext() *gin.Context {
	return uri.ctx
}

func (uri *defaultURI) Context() *gin.Context {
	return uri.ctx
}

type Can struct {
	*gin.Engine
}

func (can *Can) Run(args ...string) {
	can.Engine.Run(args...)
}

func (can *Can) Controller(uri URI) {
	typ := reflect.TypeOf(uri)
	ctrlPrefixes, ok := uriValue(typ)
	if !ok {
		return
	}
	fmt.Println(ctrlPrefixes)
	typ = reflect.PtrTo(typ)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		// 限定只有一个输入
		if method.Type.NumIn() != 2 || method.Type.NumOut() != 1 {
			continue
		}

		in := method.Type.In(1)
		parameterPrefixes, ok := uriValue(in)
		if !ok {
			continue
		}
		fmt.Println("=>", parameterPrefixes, httpMethods(reflect.New(in).Interface()))
		uris := []string{}
		for _, ctrl := range ctrlPrefixes {
			for _, param := range parameterPrefixes {
				uris = append(uris, filepath.Clean(ctrl+"/"+param))
			}
		}
		for _, u := range uris {
			for _, httpMethod := range httpMethods(reflect.New(in).Interface()) {
				can.Handle(httpMethod, u, func(ctx *gin.Context) {
					callIn := make([]reflect.Value, method.Type.NumIn())
					callIn[0] = reflect.New(typ).Elem()
					callIn[1] = reflect.New(in).Elem()
					callIn[1].FieldByName("URI").Set(reflect.ValueOf(&defaultURI{ctx: ctx}))
					psInterface := callIn[1].Addr().Interface()
					err := ctx.ShouldBind(psInterface)
					if err != nil {
						// TODO error handle
					}
					// Parameters in path
					if len(ctx.Params) > 0 {
						err = ctx.ShouldBindUri(psInterface)
						if err != nil {
							// TODO error handle
						}
					}
					outs := method.Func.Call(callIn)
					if outs[0].CanInterface() && outs[0].Interface() != nil {
						ctx.JSON(http.StatusOK, outs[0].Interface())
					}
				})
			}
		}
	}
}

func uriValue(typ reflect.Type) (prefixes []string, isCango bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		// log
		return
	}
	prefix := ""
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Implements(reflect.TypeOf((*URI)(nil)).Elem()) {
			prefix = f.Tag.Get("value")
			isCango = true
			break
		}
	}
	return strings.Split(prefix, ";"), isCango
}

func NewCan() *Can {
	return &Can{
		gin.Default(),
	}
}
