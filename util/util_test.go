/*
 * Copyright 1999-2019 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"runtime"
	"testing"
)

func TestIsExist_ForMountPoint(t *testing.T) {
	// Skip this test on Windows as it tests Unix-specific mount points
	if runtime.GOOS == "windows" {
		t.Skip("Skipping mount point test on Windows")
	}

	tests := []struct {
		device string
		want   bool
	}{
		{"/", true},
		{"/dev", true},
		{"devfs", false},
	}
	for _, tt := range tests {
		if got := IsExist(tt.device); got != tt.want {
			t.Errorf("unexpected result: %t, expected: %t", got, tt.want)
		}
	}
}
