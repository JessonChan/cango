package cango

import (
	"reflect"
	"testing"
)

func Test_appendPaths(t *testing.T) {
	type args struct {
		tagPaths []string
		strUrls  []string
		prefix   string
	}
	tests := []struct {
		name      string
		args      args
		wantPaths []string
	}{
		{
			name: "empty_tag",
			args: args{
				tagPaths: nil,
				strUrls:  []string{"str"},
				prefix:   "/prefix",
			},
			wantPaths: []string{"/prefix/str"},
		},
		{
			name: "empty_str",
			args: args{
				tagPaths: []string{"tag"},
				strUrls:  []string{},
				prefix:   "/prefix",
			},
			wantPaths: []string{"/prefix/tag"},
		},
		{
			name: "empty",
			args: args{
				tagPaths: nil,
				strUrls:  nil,
				prefix:   "/prefix",
			},
			wantPaths: []string{"/prefix"},
		},
		{
			name: "full",
			args: args{
				tagPaths: []string{"tag1", "tag2"},
				strUrls:  []string{"str1", "str2"},
				prefix:   "/prefix",
			},
			wantPaths: []string{"/prefix/str1/tag1", "/prefix/str2/tag1", "/prefix/str1/tag2", "/prefix/str2/tag2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPaths := appendPaths(tt.args.tagPaths, tt.args.strUrls, tt.args.prefix); !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("appendPaths() = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}
}
