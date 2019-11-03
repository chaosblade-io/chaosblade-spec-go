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

package channel

import (
	"context"
	"fmt"
	"testing"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

// MockLocalChannel for testing
type MockLocalChannel struct {
	Response         *spec.Response
	ScriptPath       string
	ExpectedCommands []string
	InvokeTime       int
	NoCheck          bool
	T                *testing.T
}

func (mlc *MockLocalChannel) Run(ctx context.Context, script, args string) *spec.Response {
	cmd := fmt.Sprintf("%s %s", script, args)
	if !mlc.NoCheck && mlc.ExpectedCommands[mlc.InvokeTime] != cmd {
		mlc.T.Errorf("unexpected command: %s, expected command: %s", cmd, mlc.ExpectedCommands[mlc.InvokeTime])
	}
	mlc.InvokeTime++
	return mlc.Response
}

func (mlc *MockLocalChannel) GetScriptPath() string {
	return mlc.ScriptPath
}
