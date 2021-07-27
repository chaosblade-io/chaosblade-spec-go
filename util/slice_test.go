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

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "testEmptySlice", args: args{items: []string{}}, want: []string{}},
		{name: "testDuplicatesSlice", args: args{items: []string{"1", "2", "3", "1", "3"}}, want: []string{"1", "2", "3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIntegerListToStringSlice(t *testing.T) {
	type args struct {
		flagName  string
		flagValue string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "split by comma", args: args{flagName: "local-port", flagValue: "8080,8081,8082"},
			want: []string{"8080", "8081", "8082"}},
		{name: "split by connector", args: args{flagName: "local-port", flagValue: "8080-8083"},
			want: []string{"8080", "8081", "8082", "8083"}},
		{name: "split by comma and connector", args: args{flagName: "local-port", flagValue: "7001,8080-8083"},
			want: []string{"7001", "8080", "8081", "8082", "8083"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntegerListToStringSlice(tt.args.flagName, tt.args.flagValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntegerListToStringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIntegerListToStringSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
