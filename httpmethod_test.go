package cango

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func _Test_gen_methods(t *testing.T) {
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
	fmt.Println(`func methods(ps any) (methods []string) {`)
	template = `	if _, ok := ps.($MethodMethod); ok {
		methods = append(methods, http.Method$Method)
	}`
	for _, m := range httpMethods {
		fmt.Println(strings.ReplaceAll(template, "$Method", strings.Title(strings.ToLower(m))))
	}
	fmt.Println(`	return
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
	rs := methods(ps)
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
