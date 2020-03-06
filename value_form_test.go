package cango

import (
	"fmt"
	"net/http"
	"testing"
)

type testForm struct {
	FormValue
	Name string
}

func (t *testForm) Construct(r *http.Request) {
	fmt.Println(t.Name)
	t.Name = "testValue"
}

func Test_formConstruct(t *testing.T) {
	type args struct {
		r  *http.Request
		cs FormValue
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				r: &http.Request{
					Form: map[string][]string{
						"Name": {"test"},
					},
				},
				cs: (*testForm)(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.cs = formConstruct(tt.args.r, tt.args.cs)
			if tt.args.cs.(*testForm).Name != "testValue" {
				t.Fail()
			}
		})
	}
}
