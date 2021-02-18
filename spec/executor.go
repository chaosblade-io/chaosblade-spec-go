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

package spec

import (
	"context"
	"fmt"
	"strings"
)

const (
	DestroyKey = "suid"
)

// ExpModel is the experiment data object
type ExpModel struct {
	// Target is experiment target
	Target string `json:"target,omitempty"`

	// Scope is the experiment scope
	Scope string `json:"scope,omitempty"`

	// ActionName is the experiment action FlagName, for example delay
	ActionName string `json:"action,omitempty"`

	// ActionFlags is the experiment action flags, for example time and offset
	ActionFlags map[string]string `json:"flags,omitempty"`

	// Programs
	ActionPrograms []string `json:"programs,omitempty"`

	// Categories
	ActionCategories []string `json:"categories,omitempty"`
}

// ExpExecutor defines the ExpExecutor interface
type Executor interface {
	// Name is used to identify the ExpExecutor
	Name() string

	// Exec is used to execute the experiment
	Exec(uid string, ctx context.Context, model *ExpModel) *Response

	// SetChannel
	SetChannel(channel Channel)
}

func (exp *ExpModel) GetFlags() string {
	flags := make([]string, 0)
	for k, v := range exp.ActionFlags {
		if v == "" {
			continue
		}
		flags = append(flags, fmt.Sprintf("--%s %s", k, v))
	}
	return strings.Join(flags, " ")
}

const UnknownUid = "unknown"

func SetDestroyFlag(ctx context.Context, suid string) context.Context {
	return context.WithValue(ctx, DestroyKey, suid)
}

// IsDestroy command
func IsDestroy(ctx context.Context) (string, bool) {
	suid := ctx.Value(DestroyKey)
	if suid == nil {
		return "", false
	}
	return suid.(string), true
}
