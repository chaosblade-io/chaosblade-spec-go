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
)

// Channel is an interface for command invocation
type Channel interface {

	// channel name unique
	Name() string

	// Run script with args and returns response that wraps the result
	Run(ctx context.Context, script, args string) *Response

	// GetScriptPath return the script path
	GetScriptPath() string

	// GetPidsByProcessCmdName returns the matched process other than the current process by the program command
	GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error)

	// GetPidsByProcessName returns the matched process other than the current process by the process keyword
	GetPidsByProcessName(processName string, ctx context.Context) ([]string, error)

	// GetPsArgs returns the ps command output format
	GetPsArgs(ctx context.Context) string

	// isAlpinePlatform returns true if the os version is alpine.
	// If the /etc/os-release file doesn't exist, the function returns false.
	IsAlpinePlatform(ctx context.Context) bool

	// IsAllCommandsAvailable returns nil,true if all commands exist
	IsAllCommandsAvailable(ctx context.Context, commandNames []string) (*Response, bool)

	// IsCommandAvailable returns true if the command exists
	IsCommandAvailable(ctx context.Context, commandName string) bool

	// ProcessExists returns true if the pid exists, otherwise return false.
	ProcessExists(pid string) (bool, error)

	// GetPidUser returns the process user by pid
	GetPidUser(pid string) (string, error)

	// GetPidsByLocalPorts returns the process ids using the ports
	GetPidsByLocalPorts(ctx context.Context, localPorts []string) ([]string, error)

	// GetPidsByLocalPort returns the process pid corresponding to the port
	GetPidsByLocalPort(ctx context.Context, localPort string) ([]string, error)
}
