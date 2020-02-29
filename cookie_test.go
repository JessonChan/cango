package cango

import (
	"net/http"
	"testing"
)

type MyCookie0 struct {
	Cookie
	build string
}

type MyCookie1 struct {
	Cookie
	build string
}

var build = "build"

func (m *MyCookie1) New(r *http.Request) {
	m.build = build
}

func Test_cookieConstruct(t *testing.T) {
	type args struct {
		r  *http.Request
		cs Cookie
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "0",
			args: args{
				r: new(http.Request),
				cs: &MyCookie0{
					Cookie: &emptyCookieConstructor{},
				},
			}},
		{
			name: "1",
			args: args{
				r: new(http.Request),
				cs: &MyCookie1{
					Cookie: &emptyCookieConstructor{},
					build:  "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookieConstruct(tt.args.r, tt.args.cs)
			switch tt.name {
			case "0":
				myCookie := tt.args.cs.(*MyCookie0)
				if myCookie.build != "" {
					t.Fail()
				}
			case "1":
				myCookie := tt.args.cs.(*MyCookie1)
				if myCookie.build != build {
					t.Fail()
				}
			}
		})
	}
}
