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

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

// MockLocalChannel for testing
type MockLocalChannel struct {
	ScriptPath string
	// mock function
	RunFunc                     func(ctx context.Context, script, args string) *spec.Response
	GetPidsByProcessCmdNameFunc func(processName string, ctx context.Context) ([]string, error)
	GetPidsByProcessNameFunc    func(processName string, ctx context.Context) ([]string, error)
	GetPsArgsFunc               func() string
	IsCommandAvailableFunc      func(commandName string) bool
	ProcessExistsFunc           func(pid string) (bool, error)
	GetPidUserFunc              func(pid string) (string, error)
	GetPidsByLocalPortsFunc     func(localPorts []string) ([]string, error)
	GetPidsByLocalPortFunc      func(localPort string) ([]string, error)
}

func NewMockLocalChannel() OsChannel {
	return &MockLocalChannel{
		ScriptPath:                  util.GetBinPath(),
		RunFunc:                     defaultRunFunc,
		GetPidsByProcessCmdNameFunc: defaultGetPidsByProcessCmdNameFunc,
		GetPidsByProcessNameFunc:    defaultGetPidsByProcessNameFunc,
		GetPsArgsFunc:               defaultGetPsArgsFunc,
		IsCommandAvailableFunc:      defaultIsCommandAvailableFunc,
		ProcessExistsFunc:           defaultProcessExistsFunc,
		GetPidUserFunc:              defaultGetPidUserFunc,
		GetPidsByLocalPortsFunc:     defaultGetPidsByLocalPortsFunc,
		GetPidsByLocalPortFunc:      defaultGetPidsByLocalPortFunc,
	}
}

func (mlc *MockLocalChannel) GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error) {
	return mlc.GetPidsByProcessCmdNameFunc(processName, ctx)
}

func (mlc *MockLocalChannel) GetPidsByProcessName(processName string, ctx context.Context) ([]string, error) {
	return mlc.GetPidsByProcessNameFunc(processName, ctx)
}

func (mlc *MockLocalChannel) GetPsArgs() string {
	return mlc.GetPsArgsFunc()
}

func (mlc *MockLocalChannel) isAlpinePlatform() bool {
	return false
}
func (mlc *MockLocalChannel) IsAllCommandsAvailable(commandNames []string) (*spec.Response, bool) {
	return nil, false
}

func (mlc *MockLocalChannel) IsCommandAvailable(commandName string) bool {
	return mlc.IsCommandAvailableFunc(commandName)
}

func (mlc *MockLocalChannel) ProcessExists(pid string) (bool, error) {
	return mlc.ProcessExistsFunc(pid)
}

func (mlc *MockLocalChannel) GetPidUser(pid string) (string, error) {
	return mlc.GetPidUserFunc(pid)
}

func (mlc *MockLocalChannel) GetPidsByLocalPorts(localPorts []string) ([]string, error) {
	return mlc.GetPidsByLocalPortsFunc(localPorts)
}

func (mlc *MockLocalChannel) GetPidsByLocalPort(localPort string) ([]string, error) {
	return mlc.GetPidsByLocalPortFunc(localPort)
}

func (mlc *MockLocalChannel) Run(ctx context.Context, script, args string) *spec.Response {
	return mlc.RunFunc(ctx, script, args)
}

func (mlc *MockLocalChannel) GetScriptPath() string {
	return mlc.ScriptPath
}

var defaultGetPidsByProcessCmdNameFunc = func(processName string, ctx context.Context) ([]string, error) {
	return []string{}, nil
}
var defaultGetPidsByProcessNameFunc = func(processName string, ctx context.Context) ([]string, error) {
	return []string{}, nil
}
var defaultGetPsArgsFunc = func() string {
	return "-eo user,pid,ppid,args"
}
var defaultIsCommandAvailableFunc = func(commandName string) bool {
	return false
}
var defaultProcessExistsFunc = func(pid string) (bool, error) {
	return false, nil
}
var defaultGetPidUserFunc = func(pid string) (string, error) {
	return "admin", nil
}
var defaultGetPidsByLocalPortsFunc = func(localPorts []string) ([]string, error) {
	return []string{}, nil
}
var defaultGetPidsByLocalPortFunc = func(localPort string) ([]string, error) {
	return []string{}, nil
}
var defaultRunFunc = func(ctx context.Context, script, args string) *spec.Response {
	return spec.ReturnSuccess("success")
}
