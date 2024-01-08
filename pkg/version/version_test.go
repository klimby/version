package version

import "testing"

func TestV_Compare(t *testing.T) {
	tests := []struct {
		name   string
		first  V
		second V
		want   int
	}{
		{
			name:   "invalid",
			first:  "invalid",
			second: "invalid",
			want:   0,
		},
		{
			name:   "invalid + invalid",
			first:  "alpha",
			second: "beta",
			want:   -1,
		},
		{
			name:   "1.1.0 == 1.1.0",
			first:  "1.1.0",
			second: "1.1.0",
			want:   0,
		},
		{
			name:   "1.2 > 1.1.9",
			first:  "1.2",
			second: "1.1.9",
			want:   1,
		},
		{
			name:   "1.1.9 < 1.2",
			first:  "1.1.9",
			second: "1.2",
			want:   -1,
		},
		{
			name:   "1.0.0 < 2.0.0",
			first:  "1.0.0",
			second: "2.0.0",
			want:   -1,
		},
		{
			name:   "2.0.0 < 2.1.0",
			first:  "2.0.0",
			second: "2.1.0",
			want:   -1,
		},
		{
			name:   "2.1.0 < 2.1.1",
			first:  "2.1.0",
			second: "2.1.1",
			want:   -1,
		},
		{
			name:   "v2.1.0 < 2.1.1",
			first:  "v2.1.0",
			second: "2.1.1",
			want:   -1,
		},
		{
			name:   "1.0.0-alpha < 1.0.0",
			first:  V("1.0.0-alpha"),
			second: V("1.0.0"),
			want:   -1,
		},
		{
			name:   "1.0.0-alpha < 1.0.0-beta",
			first:  V("1.0.0-alpha"),
			second: V("1.0.0-beta"),
			want:   -1,
		},
		{
			name:   "1.0.0-alpha < 1.0.0-alpha.1",
			first:  V("1.0.0-alpha"),
			second: V("1.0.0-alpha.1"),
			want:   -1,
		},
		{
			name:   "1.0.0-alpha.1 < 1.0.0-alpha.beta",
			first:  V("1.0.0-alpha.1"),
			second: V("1.0.0-alpha.beta"),
			want:   -1,
		},
		{
			name:   "1.0.0-alpha.beta < 1.0.0-beta",
			first:  V("1.0.0-alpha.beta"),
			second: V("1.0.0-beta"),
			want:   -1,
		},
		{
			name:   "1.0.0-beta < 1.0.0-beta.2",
			first:  V("1.0.0-beta"),
			second: V("1.0.0-beta.2"),
			want:   -1,
		},
		{
			name:   "1.0.0-beta.2 < 1.0.0-beta.11",
			first:  V("1.0.0-beta.2"),
			second: V("1.0.0-beta.11"),
			want:   -1,
		},
		{
			name:   "1.0.0-beta.11 < 1.0.0-rc.1",
			first:  V("1.0.0-beta.11"),
			second: V("1.0.0-rc.1"),
			want:   -1,
		},
		{
			name:   "1.0.0-rc.1 < 1.0.0",
			first:  V("1.0.0-rc.1"),
			second: V("1.0.0"),
			want:   -1,
		},
		{
			name:   "1.0.0-beta.2+gamma < 1.0.0-beta.2",
			first:  V("1.0.0-beta.2+gamma"),
			second: V("1.0.0-beta.2"),
			want:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.first.Compare(tt.second); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_Invalid(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want bool
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: true,
		},
		{
			name: "invalid",
			v:    "",
			want: true,
		},
		{
			name: "without patch",
			v:    "1.2",
			want: false,
		},
		{
			name: "without patch - alpha",
			v:    "1.2-alpha",
			want: false,
		},
		{
			name: "without patch - alpha.1",
			v:    "1.2-alpha.1",
			want: false,
		},
		{
			name: "without patch - alpha.beta",
			v:    "1.2-alpha.beta",
			want: false,
		},
		{
			name: "without patch - alpha.beta+gamma",
			v:    "1.2-alpha.beta+gamma",
			want: false,
		},

		{
			name: "v + without patch",
			v:    "v1.2",
			want: false,
		},
		{
			name: "v + without patch - alpha",
			v:    "v1.2-alpha",
			want: false,
		},
		{
			name: "v + without patch - alpha.1",
			v:    "v1.2-alpha.1",
			want: false,
		},
		{
			name: "v + without patch - alpha.beta",
			v:    "v1.2-alpha.beta",
			want: false,
		},
		{
			name: "valid",
			v:    "1.2.3",
			want: false,
		},
		{
			name: "valid - alpha",
			v:    "1.2.3-alpha",
			want: false,
		},
		{
			name: "valid - alpha.1",
			v:    "1.2.3-alpha.1",
			want: false,
		},
		{
			name: "valid - alpha.beta",
			v:    "1.2.3-alpha.beta",
			want: false,
		},
		{
			name: "valid - alpha.beta+gamma",
			v:    "1.2.3-alpha.beta+gamma",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Invalid(); got != tt.want {
				t.Errorf("Invalid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_Semver(t *testing.T) {
	tests := []struct {
		name              string
		v                 V
		wantMajor         int
		wantMinor         int
		wantPatch         int
		wantPrerelease    string
		wantBuildmetadata string
	}{
		{
			name:              "invalid",
			v:                 "invalid",
			wantMajor:         0,
			wantMinor:         0,
			wantPatch:         0,
			wantPrerelease:    "",
			wantBuildmetadata: "",
		},
		{
			name:              "without patch",
			v:                 "1.2",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "",
			wantBuildmetadata: "",
		},
		{
			name:              "without patch - alpha",
			v:                 "1.2-alpha",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha",
			wantBuildmetadata: "",
		},
		{
			name:              "without patch - alpha.1",
			v:                 "1.2-alpha.1",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.1",
			wantBuildmetadata: "",
		},
		{
			name:              "without patch - alpha.beta",
			v:                 "1.2-alpha.beta",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.beta",
			wantBuildmetadata: "",
		},
		{
			name:              "without patch - alpha.beta+gamma",
			v:                 "1.2-alpha.beta+gamma",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.beta",
			wantBuildmetadata: "gamma",
		},
		{
			name:              "v + without patch",
			v:                 "v1.2",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "",
			wantBuildmetadata: "",
		},
		{
			name:              "v + without patch - alpha",
			v:                 "v1.2-alpha",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha",
			wantBuildmetadata: "",
		},
		{
			name:              "v + without patch - alpha.1",
			v:                 "v1.2-alpha.1",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.1",
			wantBuildmetadata: "",
		},
		{
			name:              "v + without patch - alpha.1",
			v:                 "v1.2-alpha.1.2.3",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.1.2.3",
			wantBuildmetadata: "",
		},
		{
			name:              "v + without patch - alpha.beta",
			v:                 "v1.2-alpha.beta",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.beta",
			wantBuildmetadata: "",
		},
		{
			name:              "v + without patch - alpha.beta+gamma",
			v:                 "v1.2-alpha.beta+gamma",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         0,
			wantPrerelease:    "alpha.beta",
			wantBuildmetadata: "gamma",
		},
		{
			name:              "valid",
			v:                 "1.2.3",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         3,
			wantPrerelease:    "",
			wantBuildmetadata: "",
		},
		{
			name:              "valid - alpha",
			v:                 "1.2.3-alpha",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         3,
			wantPrerelease:    "alpha",
			wantBuildmetadata: "",
		},
		{
			name:              "valid - alpha.1",
			v:                 "1.2.3-alpha.1",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         3,
			wantPrerelease:    "alpha.1",
			wantBuildmetadata: "",
		},
		{
			name:              "valid - alpha.beta",
			v:                 "1.2.3-alpha.beta",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         3,
			wantPrerelease:    "alpha.beta",
			wantBuildmetadata: "",
		},
		{
			name:              "valid - alpha.beta+gamma",
			v:                 "1.2.3-alpha.beta+gamma",
			wantMajor:         1,
			wantMinor:         2,
			wantPatch:         3,
			wantPrerelease:    "alpha.beta",
			wantBuildmetadata: "gamma",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMajor, gotMinor, gotPatch, gotPrerelease, gotBuildmetadata := tt.v.semver()
			if gotMajor != tt.wantMajor {
				t.Errorf("Semver() gotMajor = %v, want %v", gotMajor, tt.wantMajor)
			}
			if gotMinor != tt.wantMinor {
				t.Errorf("Semver() gotMinor = %v, want %v", gotMinor, tt.wantMinor)
			}
			if gotPatch != tt.wantPatch {
				t.Errorf("Semver() gotPatch = %v, want %v", gotPatch, tt.wantPatch)
			}
			if gotPrerelease != tt.wantPrerelease {
				t.Errorf("Semver() gotPrerelease = %v, want %v", gotPrerelease, tt.wantPrerelease)
			}
			if gotBuildmetadata != tt.wantBuildmetadata {
				t.Errorf("Semver() gotBuildmetadata = %v, want %v", gotBuildmetadata, tt.wantBuildmetadata)
			}
		})
	}
}

func TestV_String(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want string
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: "invalid",
		},
		{
			name: "valid",
			v:    "1.0.0",
			want: "1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareASC(t *testing.T) {
	type args struct {
		v HasVersion
		o HasVersion
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1.0.0 less 2.0.0",
			args: args{
				v: V("1.0.0"),
				o: V("2.0.0"),
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareASC(tt.args.v, tt.args.o); got != tt.want {
				t.Errorf("CompareASC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareDESC(t *testing.T) {
	type args struct {
		v HasVersion
		o HasVersion
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1.0.0 greater 2.0.0",
			args: args{
				v: V("1.0.0"),
				o: V("2.0.0"),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareDESC(tt.args.v, tt.args.o); got != tt.want {
				t.Errorf("CompareDESC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_Empty(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want bool
	}{
		{
			name: "empty",
			v:    "",
			want: true,
		},
		{
			name: "not empty",
			v:    "1.0.0",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Empty(); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_Equal(t *testing.T) {
	type args struct {
		o V
	}
	tests := []struct {
		name string
		v    V
		args args
		want bool
	}{
		{
			name: "equal",
			v:    "1.0.0",
			args: args{
				o: "1.0.0",
			},
			want: true,
		},
		{
			name: "not equal",
			v:    "1.0.0",
			args: args{
				o: "2.0.0",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Equal(tt.args.o); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_FormatString(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want string
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: "",
		},
		{
			name: "valid",
			v:    "1.0.0",
			want: "1.0.0",
		},
		{
			name: "valid - v",
			v:    "v1.0.0",
			want: "1.0.0",
		},
		{
			name: "valid - alpha",
			v:    "1.0.0-alpha",
			want: "1.0.0-alpha",
		},
		{
			name: "valid - alpha.1",
			v:    "1.0.0-alpha.1",
			want: "1.0.0-alpha.1",
		},
		{
			name: "valid - alpha+beta",
			v:    "1.0.0-alpha+beta",
			want: "1.0.0-alpha+beta",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.FormatString(); got != tt.want {
				t.Errorf("FormatString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_NextMajor(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want V
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: "1.0.0",
		},
		{
			name: "valid",
			v:    "1.0.0",
			want: "2.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.NextMajor(); got != tt.want {
				t.Errorf("NextMajor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_NextMinor(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want V
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: "0.1.0",
		},
		{
			name: "valid",
			v:    "1.0.0",
			want: "1.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.NextMinor(); got != tt.want {
				t.Errorf("NextMinor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_NextPatch(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want V
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: "0.0.1",
		},
		{
			name: "valid",
			v:    "1.0.0",
			want: "1.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.NextPatch(); got != tt.want {
				t.Errorf("NextPatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_Start(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want V
	}{
		{
			name: "invalid",
			v:    "invalid",
			want: "0.0.0",
		},
		{
			name: "valid",
			v:    "1.0.0",
			want: "0.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Start(); got != tt.want {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_GitVersion(t *testing.T) {
	tests := []struct {
		name string
		v    V
		want string
	}{
		{
			name: "1.0.0",
			v:    "1.0.0",
			want: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.GitVersion(); got != tt.want {
				t.Errorf("GitVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV_LessThen(t *testing.T) {
	type args struct {
		o V
	}
	tests := []struct {
		name string
		v    V
		args args
		want bool
	}{
		{
			name: "1.0.0 < 2.0.0",
			v:    "1.0.0",
			args: args{
				o: "2.0.0",
			},
			want: true,
		},
		{
			name: "1.0.0 > 2.0.0",
			v:    "2.0.0",
			args: args{
				o: "1.0.0",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.LessThen(tt.args.o); got != tt.want {
				t.Errorf("LessThen() = %v, want %v", got, tt.want)
			}
		})
	}
}
