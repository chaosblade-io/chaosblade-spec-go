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
	"reflect"
	"testing"
)

func TestRemove(t *testing.T) {
	type args struct {
		items []string
		idx   int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: struct {
				items []string
				idx   int
			}{items: []string{"1", "2", "3"}, idx: 2},
			want: []string{"1", "2"},
		},
		{
			args: struct {
				items []string
				idx   int
			}{items: []string{"1", "2", "3"}, idx: 0},
			want: []string{"3", "2"},
		},
		{
			args: struct {
				items []string
				idx   int
			}{items: []string{"1", "2", "3"}, idx: 1},
			want: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove(tt.args.items, tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
