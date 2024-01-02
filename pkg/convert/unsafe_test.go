package convert

import (
	"reflect"
	"testing"
)

func TestBytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "bytes to string",
			args: args{
				b: []byte("test"),
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := B2S(tt.args.b); got != tt.want {
				t.Errorf("B2S() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "string to bytes",
			args: args{
				s: "test",
			},
			want: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := S2B(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S2B() = %v, want %v", got, tt.want)
			}
		})
	}
}
