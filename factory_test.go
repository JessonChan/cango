package cango

import (
	"reflect"
	"testing"
)

type Struct struct {
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
