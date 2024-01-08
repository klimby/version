package convert

import "testing"

func TestI2S(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int to string",
			args: args{
				i: 1,
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := I2S(tt.args.i); got != tt.want {
				t.Errorf("I2S() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS2Clear(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "clear string",
			args: args{
				s: "test",
			},
			want: "test",
		},
		{
			name: "clear string with tab",
			args: args{
				s: "test\t",
			},
			want: "test",
		},
		{
			name: "clear string with new line",
			args: args{
				s: "test\n",
			},
			want: "test",
		},
		{
			name: "clear string with tab and new line",
			args: args{
				s: "test\t\n",
			},
			want: "test",
		},
		{
			name: "clear string with tab and new line and space",
			args: args{
				s: "  test\t\n  test    ",
			},
			want: "test test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := S2Clear(tt.args.s); got != tt.want {
				t.Errorf("S2Clear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS2Int(t *testing.T) {
	type args struct {
		s string
		d []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "string to int",
			args: args{
				s: "1",
				d: []int{0},
			},
			want: 1,
		},
		{
			name: "string to int with default",
			args: args{
				s: "a",
				d: []int{0},
			},
			want: 0,
		},
		{
			name: "string to int without default error",
			args: args{
				s: "a",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := S2Int(tt.args.s, tt.args.d...); got != tt.want {
				t.Errorf("S2Int() = %v, want %v", got, tt.want)
			}
		})
	}
}
