package cango

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type PageController struct {
	URI       `value:"/blog/{blogName}/article/{articleId}-{pageId}"`
	BlogName  string
	ArticleId int
	PageId    int
}

type Comment struct {
	Id      int
	Content string
	Author  string
}

func (p *PageController) GetComments(param struct {
	URI    `value:"/comments/list.json"`
	Method `value:"get,post"`
	PostMethod
}) {
	fmt.Println(p.BlogName, p.ArticleId, p.PageId)
}

func (p *PageController) GetComment(param struct {
	URI       `value:"/comment/{commentId}.json"`
	CommentId int
}) Comment {
	fmt.Println("getComment", param.CommentId, p.ArticleId, p.PageId, p.BlogName)
	return Comment{Id: param.CommentId, Content: "写的太好了，五星推荐", Author: "追求理想的风"}
}

func TestAddr_String(t *testing.T) {
	type fields struct {
		Host string
		Port int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "local", fields: fields{Host: "127.0.0.1", Port: 8080}, want: "127.0.0.1:8080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := Addr{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if got := addr.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCan_Do(t *testing.T) {
	type fields struct {
		srv *http.Server
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "do", fields: fields{&http.Server{}}, args: args{&http.Request{URL: &url.URL{Path: "/blog/jack/article/501-3/comment/123.json"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			can := &Can{
				srv: tt.fields.srv,
			}
			can.Route(&PageController{})
			rt, statusCode := can.serve(tt.args.req)
			t.Log(rt, statusCode)
		})
	}
}

func TestCan_Route(t *testing.T) {
	type fields struct {
		srv *http.Server
	}
	type args struct {
		uri URI
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// can := &Can{
			// 	srv: tt.fields.srv,
			// }
		})
	}
}

func TestCan_Run(t *testing.T) {
	type fields struct {
		srv *http.Server
	}
	type args struct {
		as []Addr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// can := &Can{
			// 	srv: tt.fields.srv,
			// }
		})
	}
}

func TestCan_ServeHTTP(t *testing.T) {
	type fields struct {
		srv *http.Server
	}
	type args struct {
		rw http.ResponseWriter
		r  *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// can := &Can{
			// 	srv: tt.fields.srv,
			// }
		})
	}
}

func TestCan_urlStr(t *testing.T) {
	type fields struct {
		srv *http.Server
	}
	type args struct {
		uri interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			can := &Can{
				srv: tt.fields.srv,
			}
			got, got1 := can.urlStr(tt.args.uri)
			if got != tt.want {
				t.Errorf("urlStr() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("urlStr() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewCan(t *testing.T) {
	tests := []struct {
		name string
		want *Can
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCan(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assignValue(t *testing.T) {
	type args struct {
		value reflect.Value
		vars  map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_context_Request(t *testing.T) {
	type fields struct {
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		want   *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &context{
				request: tt.fields.request,
			}
			if got := c.Request(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lowerCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Lower", args: args{"Lower"}, want: "lower"},
		{name: "", args: args{""}, want: ""},
		{name: "Big", args: args{"Big"}, want: "big"},
		{name: "big", args: args{"big"}, want: "big"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowerCase(tt.args.str); got != tt.want {
				t.Errorf("lowerCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name     string
		args     args
		wantVars map[int]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVars := parseTag(tt.args.tag); !reflect.DeepEqual(gotVars, tt.wantVars) {
				t.Errorf("parseTag() = %v, want %v", gotVars, tt.wantVars)
			}
		})
	}
}

func Test_upperCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "upper", args: args{"upper"}, want: "Upper"},
		{name: "", args: args{""}, want: ""},
		{name: "camel", args: args{"camel"}, want: "Camel"},
		{name: "Camel", args: args{"Camel"}, want: "Camel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upperCase(tt.args.str); got != tt.want {
				t.Errorf("upperCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
