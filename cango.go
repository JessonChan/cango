package cango

import (
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type URI interface {
	ginContext() *gin.Context
	Context() *gin.Context
}

// TODO MustBind
var DefaultBind = func(ctx *gin.Context, i any) error {
	return ctx.ShouldBind(i)
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

func (can *Can) Run(args ...string) error {
	return can.Engine.Run(append(args, ":"+ternary(viper.GetString("default.port") != "", viper.GetString("default.port"), viper.GetString("cango.port")))[0])
}

func (can *Can) Controller(uri URI) {
	typ := reflect.TypeOf(uri)
	ctrlPrefixes, ok := uriValue(typ)
	if !ok {
		return
	}
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
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
		urls := []string{}
		for _, ctrl := range ctrlPrefixes {
			for _, param := range parameterPrefixes {
				urls = append(urls, filepath.Clean(ctrl+"/"+param))
			}
		}
		for _, url := range urls {
			for _, httpMethod := range httpMethods(reflect.New(in).Interface()) {
				can.Handle(httpMethod, url, func(ctx *gin.Context) {
					callIn := make([]reflect.Value, method.Type.NumIn())
					callIn[0] = reflect.New(typ).Elem()
					callIn[1] = reflect.New(in).Elem()
					callIn[1].FieldByName("URI").Set(reflect.ValueOf(&defaultURI{ctx: ctx}))
					psInterface := callIn[1].Addr().Interface()
					bind := (ShouldBind{}).bind()
					if v, ok := uri.(requestBind); ok {
						bind = v.bind()
					}
					err := bind(ctx, psInterface)
					if err != nil {
						// TODO error handle
					}
					// Parameters in path
					if len(ctx.Params) > 0 {
						err = bind(ctx, psInterface)
						if err != nil {
							// TODO error handle
						}
					}
					outs := method.Func.Call(callIn)
					if len(outs) > 0 && outs[0].CanInterface() && outs[0].Interface() != nil {
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

func NewCan(args ...string) *Can {
	//TODO viper args
	viper.SetConfigName("cango")
	viper.SetConfigType("ini")
	viper.AddConfigPath("./conf/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}
	return &Can{
		gin.Default(),
	}
}

// ternary for example ternary(true, "a", "b")=="a"
func ternary[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}
