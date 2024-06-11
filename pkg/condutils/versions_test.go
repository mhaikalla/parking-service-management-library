package condutils

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewSemanticVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    SemanticVersion
		wantErr bool
	}{
		{"1", args{version: "0.0.1"}, SemanticVersion{major: 0, minor: 0, patch: 1, versionLong: 1}, false},
		{"2", args{version: "0.1.1"}, SemanticVersion{major: 0, minor: 1, patch: 1, versionLong: 11}, false},
		{"3", args{version: "1.1.1"}, SemanticVersion{major: 1, minor: 1, patch: 1, versionLong: 111}, false},
		{"4", args{version: "01.001.0001"}, SemanticVersion{major: 1, minor: 1, patch: 1, versionLong: 10010001}, false},
		{"5", args{version: "1.1.1A"}, SemanticVersion{}, true},
		{"6", args{version: "1.A.1"}, SemanticVersion{}, true},
		{"7", args{version: "1A.1.10"}, SemanticVersion{}, true},
		{"8", args{version: "1.1"}, SemanticVersion{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSemanticVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSemanticVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSemanticVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionAfterComparation(t *testing.T) {
	type args struct {
		current string
		test    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{current: "0.0.2", test: "0.0.1"}, true},
		{"2", args{current: "0.0.1", test: "0.0.2"}, false},
		{"3", args{current: "0.0.1", test: "0.0.1"}, false},
		{"4", args{current: "0.0.0", test: "0.0.0"}, false},

		{"5", args{current: "0.2.0", test: "0.1.0"}, true},
		{"6", args{current: "0.1.0", test: "0.2.0"}, false},
		{"7", args{current: "0.1.0", test: "0.1.0"}, false},
		{"8", args{current: "0.1.10", test: "0.02.1"}, false},
		{"9", args{current: "0.005.0", test: "0.02.1"}, true},

		{"10", args{current: "2.0.0", test: "1.0.0"}, true},
		{"11", args{current: "1.0.0", test: "2.0.0"}, false},
		{"12", args{current: "1.0.0", test: "1.0.0"}, false},
		{"13", args{current: "1.1.0", test: "0.0002.1"}, true},
		{"14", args{current: "0003.005.0", test: "0001.02.10000"}, true},
		{"15", args{current: "0001.005.000001", test: "0001.00010.10000"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, err1 := NewSemanticVersion(tt.args.current)
			if err1 != nil {
				t.Errorf("NewSemanticVersion() error = %v", err1)
				return
			}

			test, err2 := NewSemanticVersion(tt.args.test)
			if err2 != nil {
				t.Errorf("NewSemanticVersion() error = %v", err2)
				return
			}

			if current.After(test) != tt.want {
				t.Log(current.After(test))
				t.Errorf("NewSemanticVersion() = current %v, test %v, current not after test", current, test)
			}
		})
	}
}

func TestVersionBeforeComparation(t *testing.T) {
	type args struct {
		current string
		test    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{current: "0.0.2", test: "0.0.1"}, false},
		{"2", args{current: "0.0.1", test: "0.0.2"}, true},
		{"3", args{current: "0.0.1", test: "0.0.1"}, false},
		{"4", args{current: "0.0.0", test: "0.0.0"}, false},

		{"5", args{current: "0.2.0", test: "0.1.0"}, false},
		{"6", args{current: "0.1.0", test: "0.2.0"}, true},
		{"7", args{current: "0.1.0", test: "0.1.0"}, false},
		{"8", args{current: "0.1.10", test: "0.02.1"}, true},
		{"9", args{current: "0.005.0", test: "0.02.1"}, false},

		{"10", args{current: "2.0.0", test: "1.0.0"}, false},
		{"11", args{current: "1.0.0", test: "2.0.0"}, true},
		{"12", args{current: "1.0.0", test: "1.0.0"}, false},
		{"13", args{current: "1.1.0", test: "0.0002.1"}, false},
		{"14", args{current: "0003.005.0", test: "0001.02.10000"}, false},
		{"15", args{current: "0001.005.000001", test: "0001.00010.10000"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, err1 := NewSemanticVersion(tt.args.current)
			if err1 != nil {
				t.Errorf("NewSemanticVersion() error = %v", err1)
				return
			}

			test, err2 := NewSemanticVersion(tt.args.test)
			if err2 != nil {
				t.Errorf("NewSemanticVersion() error = %v", err2)
				return
			}

			if current.Before(test) != tt.want {
				t.Log(current.After(test))
				t.Errorf("NewSemanticVersion() = current %v, test %v, current not after test", current, test)
			}
		})
	}
}

func TestAllSemanticVersionAllMethod(t *testing.T) {
	ver, parseErr := NewSemanticVersion("5.1.0")
	if parseErr != nil {
		t.FailNow()
	}

	if ver.Major() != 5 {
		t.FailNow()
	}

	if ver.Minor() != 1 {
		t.FailNow()
	}

	if ver.Patch() != 0 {
		t.FailNow()
	}
	if ver.ToString() == "" ||
		!regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(ver.ToString()) {
		t.FailNow()
	}
}
