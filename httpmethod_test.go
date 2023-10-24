package cango

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func test_gen_methods(t *testing.T) {
	httpMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
		http.MethodTrace,
	}
	template := `type $MethodMethod interface {
		cangoHttp$Method()
	}
	`
	for _, m := range httpMethods {
		fmt.Println(strings.ReplaceAll(template, "$Method", strings.Title(strings.ToLower(m))))
	}
	fmt.Println(`func httpMethods(ps any) (methods []string) {`)
	template = `	if _, ok := ps.($MethodMethod); ok {
		methods = append(methods, http.Method$Method)
	}`
	for _, m := range httpMethods {
		fmt.Println(strings.ReplaceAll(template, "$Method", strings.Title(strings.ToLower(m))))
	}
	fmt.Println(`		// 默认为Get
	if len(methods) == 0 {
		methods = append(methods, http.MethodGet)
	}
	return
	}`)
}
func Test_Methods(t *testing.T) {
	ps := struct {
		GetMethod
		PostMethod
		PatchMethod
		HeadMethod
		DeleteMethod
		PutMethod
		OptionsMethod
		TraceMethod
	}{}
	rs := httpMethods(ps)
	should := strings.Split("GET POST PUT DELETE PATCH OPTIONS HEAD TRACE", " ")
	if len(rs) != len(should) {
		t.Fail()
	}
	for k, v := range should {
		if rs[k] != v {
			t.Fail()
		}
	}
}
