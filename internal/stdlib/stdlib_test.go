// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdlib

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/discovery/internal/thirdparty/semver"
)

func TestTagForVersion(t *testing.T) {
	for _, tc := range []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "std version v1.12.5",
			version: "v1.12.5",
			want:    "go1.12.5",
		},
		{
			name:    "std version v1.13, incomplete canonical version",
			version: "v1.13",
			want:    "go1.13",
		},
		{
			name:    "std version v1.13.0-beta.1",
			version: "v1.13.0-beta.1",
			want:    "go1.13beta1",
		},
		{
			name:    "std with digitless prerelease",
			version: "v1.13.0-prerelease",
			want:    "go1.13prerelease",
		},
		{
			name:    "version v1.13.0",
			version: "v1.13.0",
			want:    "go1.13",
		},
		{
			name:    "bad std semver",
			version: "v1.x",
			wantErr: true,
		},
		{
			name:    "more bad std semver",
			version: "v1.0-",
			wantErr: true,
		},
		{
			name:    "bad prerelease",
			version: "v1.13.0-beta1",
			wantErr: true,
		},
		{
			name:    "another bad prerelease",
			version: "v1.13.0-whatevs99",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := TagForVersion(tc.version)
			if (err != nil) != tc.wantErr {
				t.Errorf("TagForVersion(%q) = %q, %v, wantErr %v", tc.version, got, err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("TagForVersion(%q) = %q, %v, wanted %q, %v", tc.version, got, err, tc.want, nil)
			}
		})
	}
}

func TestZip(t *testing.T) {
	UseTestData = true
	defer func() { UseTestData = false }()

	for _, version := range []string{"v1.12.5", "v1.3.2"} {
		t.Run(version, func(t *testing.T) {
			zr, gotTime, err := Zip(version)
			if err != nil {
				t.Fatal(err)
			}
			if !gotTime.Equal(TestCommitTime) {
				t.Errorf("commit time: got %s, want %s", gotTime, TestCommitTime)
			}
			wantFiles := map[string]bool{
				"LICENSE":               true,
				"errors/errors.go":      true,
				"errors/errors_test.go": true,
			}
			if semver.Compare(version, "v1.4.0") > 0 {
				wantFiles["README.md"] = true
			} else {
				wantFiles["README"] = true
			}

			wantPrefix := "std@" + version + "/"
			for _, f := range zr.File {
				if !strings.HasPrefix(f.Name, wantPrefix) {
					t.Errorf("filename %q missing prefix %q", f.Name, wantPrefix)
					continue
				}
				delete(wantFiles, f.Name[len(wantPrefix):])
			}
			if len(wantFiles) > 0 {
				t.Errorf("zip missing files: %v", reflect.ValueOf(wantFiles).MapKeys())
			}
		})
	}
}

func TestReleaseVersionForTag(t *testing.T) {
	for _, tc := range []struct {
		in, want string
	}{
		{"", ""},
		{"go1.9beta2", ""},
		{"go1.12", "v1.12.0"},
		{"go1.9.7", "v1.9.7"},
		{"go2.0", "v2.0.0"},
	} {
		got := releaseVersionForTag(tc.in)
		if got != tc.want {
			t.Errorf("releaseVersionForTag(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}