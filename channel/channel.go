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
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

// Run the script in the local with ctx and returns spec.Response
func Run(ctx context.Context, script, args string) *spec.Response {
	return execScript(ctx, script, args)
}

// GetScriptPath returns the script path for running
func GetScriptPath() string {
	return util.GetBinPath()
}

// execScript invokes exec.CommandContext
func execScript(ctx context.Context, script, args string) *spec.Response {
	newCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if ctx == context.Background() {
		ctx = newCtx
	}
	script = strings.Replace(script, " ", `\ `, -1)
	logrus.Debugf("script: %s %s", script, args)
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script+" "+args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf(string(output) + " " + err.Error())
		return spec.ReturnFail(spec.Code[spec.ExecCommandError], errMsg)
	}
	result := string(output)
	return spec.ReturnSuccess(result)
}

// GetPidsByProcessCmdName returns the matched process other than the current process
func GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error) {
	excludeProcess := ctx.Value(ExcludeProcessKey)
	excludeGrepInfo := ""
	if excludeProcess != nil {
		excludeProcessString := excludeProcess.(string)
		if excludeProcessString != "" {
			excludeGrepInfo = fmt.Sprintf(`| grep -v -w %s`, excludeProcessString)
		}
	}
	response := Run(ctx, "pgrep",
		fmt.Sprintf(`-l %s %s | grep -v -w chaos_killprocess | grep -v -w chaos_stopprocess | awk '{print $1}' | tr '\n' ' '`,
			processName, excludeGrepInfo))
	if !response.Success {
		return nil, fmt.Errorf(response.Err)
	}
	pidString := response.Result.(string)
	pids := strings.Fields(strings.TrimSpace(pidString))
	currPid := strconv.Itoa(os.Getpid())
	for idx, pid := range pids {
		if pid == currPid {
			return util.Remove(pids, idx), nil
		}
	}
	return pids, nil
}

// grep ${key}
const ProcessKey = "process"
const ExcludeProcessKey = "excludeProcess"

// GetPidsByProcessName returns the matched process other than the current process
func GetPidsByProcessName(processName string, ctx context.Context) ([]string, error) {
	psArgs := GetPsArgs()
	otherProcess := ctx.Value(ProcessKey)
	otherGrepInfo := ""
	if otherProcess != nil {
		processString := otherProcess.(string)
		if processString != "" {
			otherGrepInfo = fmt.Sprintf(`| grep "%s"`, processString)
		}
	}
	excludeProcess := ctx.Value(ExcludeProcessKey)
	excludeGrepInfo := ""
	if excludeProcess != nil {
		excludeProcessString := excludeProcess.(string)
		if excludeProcessString != "" {
			excludeGrepInfo = fmt.Sprintf(`| grep -v -w %s`, excludeProcessString)
		}
	}
	response := Run(ctx, "ps",
		fmt.Sprintf(`%s | grep "%s" %s %s | grep -v -w grep | grep -v -w chaos_killprocess | grep -v -w chaos_stopprocess | awk '{print $2}' | tr '\n' ' '`,
			psArgs, processName, otherGrepInfo, excludeGrepInfo))
	if !response.Success {
		return nil, fmt.Errorf(response.Err)
	}
	pidString := strings.TrimSpace(response.Result.(string))
	if pidString == "" {
		return make([]string, 0), nil
	}
	pids := strings.Fields(pidString)
	currPid := strconv.Itoa(os.Getpid())
	for idx, pid := range pids {
		if pid == currPid {
			return util.Remove(pids, idx), nil
		}
	}
	return pids, nil
}

// GetPsArgs for querying the process info
func GetPsArgs() string {
	var psArgs = "-eo user,pid,ppid,args"
	if isAlpinePlatform() {
		psArgs = "-o user,pid,ppid,args"
	}
	return psArgs
}

// isAlpinePlatform returns true if the os version is alpine.
// If the /etc/os-release file doesn't exist, the function returns false.
func isAlpinePlatform() bool {
	var osVer = ""
	if util.IsExist("/etc/os-release") {
		response := Run(context.TODO(), "awk", "-F '=' '{if ($1 == \"ID\") {print $2;exit 0}}' /etc/os-release")
		if response.Success {
			osVer = response.Result.(string)
		}
	}
	return strings.TrimSpace(osVer) == "alpine"
}

// IsCommandAvailable return true if the command exists
func IsCommandAvailable(commandName string) bool {
	response := execScript(context.TODO(), "command", fmt.Sprintf("-v %s", commandName))
	return response.Success
}

//ProcessExists returns true if the pid exists, otherwise return false.
func ProcessExists(pid string) (bool, error) {
	if isAlpinePlatform() {
		response := Run(context.TODO(), "ps", fmt.Sprintf("-o pid | grep %s", pid))
		if !response.Success {
			return false, fmt.Errorf(response.Err)
		}
		if strings.TrimSpace(response.Result.(string)) == "" {
			return false, nil
		}
		return true, nil
	}
	response := Run(context.TODO(), "ps", fmt.Sprintf("-p %s", pid))
	return response.Success, nil
}

// GetPidUser returns the process user by pid
func GetPidUser(pid string) (string, error) {
	var response *spec.Response
	if isAlpinePlatform() {
		response = Run(context.TODO(), "ps", fmt.Sprintf("-o user,pid | grep %s", pid))

	} else {
		response = Run(context.TODO(), "ps", fmt.Sprintf("-o user,pid -p %s | grep %s", pid, pid))
	}
	if !response.Success {
		return "", fmt.Errorf(response.Err)
	}
	result := strings.TrimSpace(response.Result.(string))
	if result == "" {
		return "", fmt.Errorf("process user not found by pid")
	}
	return strings.Fields(result)[0], nil
}
