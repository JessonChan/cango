package cango

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_parseString(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *fastMux
	}{{
		args: args{"/Users/jessonchan////code/go/path/src/github.com/JessonChan/cango/examples/can/run.go"},
	}, {
		args: args{"weather/2020-02-01/how-heavy/very/rain.json"},
	}, {
		args: args{"/weather/2020-02-01/how-heavy/very/snowing/white.json"},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePath(tt.args.url)
			t.Log(got)
		})
	}
}

func Test_parseString_01(t *testing.T) {
	mux := newMux()
	router := mux.newRouter("snow")
	router.path("/weather/{day}/././how-heavy/{heavy}/snowing/{color}.json/")
	router.methods(http.MethodGet)

	mux.doMatch("GET", "/weather/2020-02-01/how-heavy/very/snowing/white.json")

	pair := mux.match(http.MethodGet, "/weather/2020-02-01/how-heavy/very/snowing/white.json")
	if pair == nil {
		t.Fatal("failed")
	}
	fmt.Println(pair.vars)
}
