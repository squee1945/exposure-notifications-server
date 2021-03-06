// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiconfig

import (
	"testing"
	"time"

	"github.com/google/exposure-notifications-server/internal/android"

	"github.com/google/go-cmp/cmp"
)

func TestBaseAPIConfig(t *testing.T) {

	cfg := New()
	if cfg.IsIOS() {
		t.Errorf("cfg.IoIOS, got true, want false")
	}
	if cfg.IsAndroid() {
		t.Errorf("cfg.IoAndroid, got true, want false")
	}
}

func TestVerifyOpts(t *testing.T) {
	testTime := time.Date(2020, 1, 13, 5, 6, 4, 6, time.Local)

	cases := []struct {
		cfg  *APIConfig
		opts android.VerifyOpts
	}{
		{
			cfg: &APIConfig{
				AppPackageName:    "foo",
				CTSProfileMatch:   true,
				BasicIntegrity:    true,
				AllowedPastTime:   time.Duration(15 * time.Minute),
				AllowedFutureTime: time.Duration(1 * time.Second),
			},
			opts: android.VerifyOpts{
				AppPkgName:      "foo",
				APKDigest:       []string{},
				CTSProfileMatch: true,
				BasicIntegrity:  true,
				MinValidTime:    testTime.Add(-15 * time.Minute),
				MaxValidTime:    testTime.Add(1 * time.Second),
			},
		},
		{
			cfg: &APIConfig{
				AppPackageName:    "foo",
				CTSProfileMatch:   false,
				BasicIntegrity:    true,
				AllowedPastTime:   0,
				AllowedFutureTime: 0,
			},
			opts: android.VerifyOpts{
				AppPkgName:      "foo",
				APKDigest:       []string{},
				CTSProfileMatch: false,
				BasicIntegrity:  true,
				MinValidTime:    time.Time{},
				MaxValidTime:    time.Time{},
			},
		},
		{
			cfg: &APIConfig{
				AppPackageName:    "foo",
				ApkDigestSHA256:   []string{"bar"},
				CTSProfileMatch:   false,
				BasicIntegrity:    true,
				AllowedPastTime:   0,
				AllowedFutureTime: 0,
			},
			opts: android.VerifyOpts{
				AppPkgName:      "foo",
				APKDigest:       []string{"bar"},
				CTSProfileMatch: false,
				BasicIntegrity:  true,
				MinValidTime:    time.Time{},
				MaxValidTime:    time.Time{},
			},
		},
	}

	for i, tst := range cases {
		got := tst.cfg.VerifyOpts(testTime, nil /* noncer */)
		if diff := cmp.Diff(tst.opts, got); diff != "" {
			t.Errorf("%v verify opts (-want +got):\n%v", i, diff)
		}
	}
}
