package fsys

import (
	"testing"

	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/viper"
)

func TestFile_Empty(t *testing.T) {
	tests := []struct {
		name string
		f    File
		want bool
	}{
		{
			name: "empty",
			f:    "",
			want: true,
		},
		{
			name: "not empty",
			f:    "test",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Empty(); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_IsAbs(t *testing.T) {
	tests := []struct {
		name string
		f    File
		want bool
	}{
		{
			name: "is abs",
			f:    "/test",
			want: true,
		},
		{
			name: "is not abs",
			f:    "test",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsAbs(); got != tt.want {
				t.Errorf("IsAbs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Path(t *testing.T) {
	tests := []struct {
		name    string
		f       File
		workDir string
		want    string
	}{
		{
			name:    "path",
			f:       "test",
			workDir: "/foo",
			want:    "/foo/test",
		},
		{
			name:    "path",
			f:       "/test",
			workDir: "/foo",
			want:    "/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.WorkDir, tt.workDir)

			if got := tt.f.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Rel(t *testing.T) {
	tests := []struct {
		name    string
		f       File
		workDir string
		want    string
	}{
		{
			name:    "rel",
			f:       "test",
			workDir: "/foo",
			want:    "test",
		},
		{
			name:    "rel",
			f:       "/test",
			workDir: "/foo",
			want:    "../test",
		},
		{
			name:    "rel",
			f:       "/test",
			workDir: "foo",
			want:    "/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.WorkDir, tt.workDir)

			if got := tt.f.Rel(); got != tt.want {
				t.Errorf("Rel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_String(t *testing.T) {
	tests := []struct {
		name string
		f    File
		want string
	}{
		{
			name: "string",
			f:    "test",
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
