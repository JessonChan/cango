package cango

import (
	"reflect"
	"testing"
)

type Struct struct {
	URI
}

type NotURI struct {
}

func Test_toPtrKind(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "struct_01",
			args: args{i: &Struct{}},
			want: reflect.TypeOf(&Struct{}),
		},
		{
			name: "struct_02",
			args: args{i: Struct{}},
			want: reflect.TypeOf(&Struct{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toPtrKind(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toPtrKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_implements(t *testing.T) {
	type args struct {
		in  reflect.Type
		typ reflect.Type
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "struct",
			args: args{
				in:  reflect.TypeOf(Struct{}),
				typ: uriType,
			},
			want: reflect.TypeOf(Struct{}),
		},
		{
			name: "ptr",
			args: args{
				in:  reflect.TypeOf(&Struct{}),
				typ: uriType,
			},
			want: reflect.TypeOf(Struct{}),
		},
		{
			name: "not_uri",
			args: args{
				in:  reflect.TypeOf(&NotURI{}),
				typ: uriType,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := implements(tt.args.in, tt.args.typ)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("implements() got = %v, want %v", got, tt.want)
			}
		})
	}
}
