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
	"fmt"
	"strconv"
	"strings"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

// Remove the item by the index
func Remove(items []string, idx int) []string {
	length := len(items)
	items[length-1], items[idx] = items[idx], items[length-1]
	return items[:length-1]
}

func RemoveDuplicates(items []string) []string {
	result := make([]string, 0)
	cache := map[string]struct{}{}
	for _, i := range items {
		if _, ok := cache[i]; !ok {
			cache[i] = struct{}{}
			result = append(result, i)
		}
	}
	return result
}

// ParseIntegerListToStringSlice func parses the multiple integer values to string slice.
// Support the below formats: 0 | 0,1 | 0,2,3 | 0-3 | 0,2-4 | 0,1,3-5
// For example, the flag value is 0,2-3, the func returns []string{"0", "2", "3"}
func ParseIntegerListToStringSlice(flagName, flagValue string) ([]string, error) {
	values := make([]string, 0)
	commaParts := strings.Split(flagValue, ",")
	for _, part := range commaParts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		if !strings.Contains(value, "-") {
			_, err := strconv.Atoi(value)
			if err != nil {
				return values, fmt.Errorf(spec.ParameterIllegal.Sprintf(flagName, flagValue, err))
			}
			values = append(values, value)
			continue
		}
		ranges := strings.Split(value, "-")
		if len(ranges) != 2 {
			return values, fmt.Errorf(spec.ParameterIllegal.Sprintf(flagName, flagValue,
				"Does not conform to the data format, a connector is required"))
		}
		startIndex, err := strconv.Atoi(strings.TrimSpace(ranges[0]))
		if err != nil {
			return values, fmt.Errorf(spec.ParameterIllegal.Sprintf(flagName, flagValue, err))
		}
		endIndex, err := strconv.Atoi(strings.TrimSpace(ranges[1]))
		if err != nil {
			return values, fmt.Errorf(spec.ParameterIllegal.Sprintf(flagName, flagValue, err))
		}
		for i := startIndex; i <= endIndex; i++ {
			values = append(values, strconv.Itoa(i))
		}
	}
	return values, nil
}
