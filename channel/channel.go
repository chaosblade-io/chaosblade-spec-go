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
)

type OsChannel interface {
	spec.Channel

	// GetPidsByProcessCmdName returns the matched process other than the current process by the program command
	GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error)

	// GetPidsByProcessName returns the matched process other than the current process by the process keyword
	GetPidsByProcessName(processName string, ctx context.Context) ([]string, error)

	// GetPsArgs returns the ps command output format
	GetPsArgs() string

	// isAlpinePlatform returns true if the os version is alpine.
	// If the /etc/os-release file doesn't exist, the function returns false.
	isAlpinePlatform() bool

	// IsAllCommandsAvailable returns nil,true if all commands exist
	IsAllCommandsAvailable(commandNames []string) (*spec.Response, bool)

	// IsCommandAvailable returns true if the command exists
	IsCommandAvailable(commandName string) bool

	// ProcessExists returns true if the pid exists, otherwise return false.
	ProcessExists(pid string) (bool, error)

	// GetPidUser returns the process user by pid
	GetPidUser(pid string) (string, error)

	// GetPidsByLocalPorts returns the process ids using the ports
	GetPidsByLocalPorts(localPorts []string) ([]string, error)

	// GetPidsByLocalPort returns the process pid corresponding to the port
	GetPidsByLocalPort(localPort string) ([]string, error)
}

// grep ${key}
const ProcessKey = "process"
const ExcludeProcessKey = "excludeProcess"
