package console

import (
	"bytes"
	"testing"
)

func TestError(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not silent not colorize",
			args: args{
				s: "test",
			},
			want: "test\n",
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "not silent colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(err) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := InitTest(tt.args.colorize)
			defer d()

			Error(tt.args.s)

			stdout, stderr := Read()

			if got := stderr; got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stdout; got != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(warning) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := InitTest(tt.args.colorize)
			defer d()

			Warn(tt.args.s)

			stdout, _ := Read()

			if got := stdout; got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(success) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := InitTest(tt.args.colorize)
			defer d()

			Success(tt.args.s)

			stdout, _ := Read()

			if got := stdout; got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestNotice(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(notice) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := InitTest(tt.args.colorize)
			defer d()

			Notice(tt.args.s)

			stdout, _ := Read()

			if got := stdout; got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(info) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := InitTest(tt.args.colorize)
			defer d()

			Info(tt.args.s)

			stdout, _ := Read()

			if got := stdout; got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "write to nowhere",
			args: args{
				p: []byte("test"),
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nilWriter{}
			got, err := n.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testWriter_Write(t *testing.T) {
	type fields struct {
		buffer bytes.Buffer
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "write to buffer",
			fields: fields{
				buffer: bytes.Buffer{},
			},
			args: args{
				p: []byte("test"),
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &testWriter{
				buffer: tt.fields.buffer,
			}
			got, err := n.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testWriter_String(t *testing.T) {
	tests := []struct {
		name  string
		write string
		want  string
	}{
		{
			name:  "return buffer as string",
			write: "test",
			want:  "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &testWriter{}

			_, err := n.Write([]byte(tt.write))
			if err != nil {
				t.Errorf("Write() error = %v", err)
				return
			}

			if got := n.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}

			// reset buffer test
			if n.buffer.Len() != 0 {
				t.Errorf("String() buffer not reset")
			}
		})
	}
}
