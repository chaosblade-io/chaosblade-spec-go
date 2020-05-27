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
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type LocalChannel struct {
}

// NewLocalChannel returns a local channel for invoking the host command
func NewLocalChannel() OsChannel {
	return &LocalChannel{}
}

func (l *LocalChannel) Run(ctx context.Context, script, args string) *spec.Response {
	return execScript(ctx, script, args)
}

func (l *LocalChannel) GetScriptPath() string {
	return util.GetBinPath()
}

func (l *LocalChannel) GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error) {
	excludeProcess := ctx.Value(ExcludeProcessKey)
	excludeGrepInfo := ""
	if excludeProcess != nil {
		excludeProcessString := excludeProcess.(string)
		if excludeProcessString != "" {
			excludeGrepInfo = fmt.Sprintf(`| grep -v -w %s`, excludeProcessString)
		}
	}
	response := l.Run(ctx, "pgrep",
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

func (l *LocalChannel) GetPidsByProcessName(processName string, ctx context.Context) ([]string, error) {
	psArgs := l.GetPsArgs()
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
	if strings.HasPrefix(processName, "-") {
		processName = fmt.Sprintf(`\%s`, processName)
	}
	response := l.Run(ctx, "ps",
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

func (l *LocalChannel) GetPsArgs() string {
	var psArgs = "-eo user,pid,ppid,args"
	if l.isAlpinePlatform() {
		psArgs = "-o user,pid,ppid,args"
	}
	return psArgs
}

func (l *LocalChannel) isAlpinePlatform() bool {
	var osVer = ""
	if util.IsExist("/etc/os-release") {
		response := l.Run(context.TODO(), "awk", "-F '=' '{if ($1 == \"ID\") {print $2;exit 0}}' /etc/os-release")
		if response.Success {
			osVer = response.Result.(string)
		}
	}
	return strings.TrimSpace(osVer) == "alpine"
}

func (l *LocalChannel) IsCommandAvailable(commandName string) bool {
	response := l.Run(context.TODO(), "command", fmt.Sprintf("-v %s", commandName))
	return response.Success
}

func (l *LocalChannel) ProcessExists(pid string) (bool, error) {
	if l.isAlpinePlatform() {
		response := l.Run(context.TODO(), "ps", fmt.Sprintf("-o pid | grep %s", pid))
		if !response.Success {
			return false, fmt.Errorf(response.Err)
		}
		if strings.TrimSpace(response.Result.(string)) == "" {
			return false, nil
		}
		return true, nil
	}
	response := l.Run(context.TODO(), "ps", fmt.Sprintf("-p %s", pid))
	return response.Success, nil
}

func (l *LocalChannel) GetPidUser(pid string) (string, error) {
	var response *spec.Response
	if l.isAlpinePlatform() {
		response = l.Run(context.TODO(), "ps", fmt.Sprintf("-o user,pid | grep %s", pid))

	} else {
		response = l.Run(context.TODO(), "ps", fmt.Sprintf("-o user,pid -p %s | grep %s", pid, pid))
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

func (l *LocalChannel) GetPidsByLocalPorts(localPorts []string) ([]string, error) {
	if localPorts == nil || len(localPorts) == 0 {
		return nil, fmt.Errorf("the local port parameter is empty")
	}
	var result = make([]string, 0)
	for _, port := range localPorts {
		pids, err := l.GetPidsByLocalPort(port)
		if err != nil {
			return nil, fmt.Errorf("failed to get pid by %s, %v", port, err)
		}
		logrus.Infof("get pids by %s port returns %v", port, pids)
		if pids != nil && len(pids) > 0 {
			result = append(result, pids...)
		}
	}
	return result, nil
}

func (l *LocalChannel) GetPidsByLocalPort(localPort string) ([]string, error) {
	available := l.IsCommandAvailable("ss")
	if !available {
		return nil, fmt.Errorf("ss command not found, can't get pid by port")
	}
	//$ss -lpn 'sport = :80'
	//Netid State      Recv-Q Send-Q   Local Address:Port   Peer Address:Port
	//tcp   LISTEN     0      128       *:80                 *:* users:(("tengine",pid=237768,fd=6),("tengine",pid=237767,fd=6))
	response := l.Run(context.TODO(), "ss", fmt.Sprintf("-pln sport = %s", localPort))
	if !response.Success {
		return []string{}, fmt.Errorf(response.Err)
	}
	if util.IsNil(response.Result) {
		return []string{}, nil
	}
	result := response.Result.(string)
	ssMsg := strings.TrimSpace(result)
	if ssMsg == "" {
		return []string{}, nil
	}
	sockets := strings.Split(ssMsg, "\n")
	logrus.Infof("sockets for %s, %v", localPort, sockets)
	for idx, s := range sockets {
		if idx == 0 {
			continue
		}
		fields := strings.Fields(s)
		// users:(("tengine",pid=237768,fd=6),("tengine",pid=237767,fd=6))
		lastField := fields[len(fields)-1]
		pidExp := regexp.MustCompile(`pid=(\d+)`)
		values := pidExp.FindStringSubmatch(lastField)
		if values == nil {
			return []string{}, nil
		}
		return values, nil
	}
	return []string{}, nil
}

// execScript invokes exec.CommandContext
func execScript(ctx context.Context, script, args string) *spec.Response {
	newCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if ctx == context.Background() {
		ctx = newCtx
	}
	isBladeCmd := isBladeCommand(script)
	script = strings.Replace(script, " ", `\ `, -1)
	logrus.Debugf("script: %s %s", script, args)
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script+" "+args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := string(output)
		if !isBladeCmd {
			errMsg = fmt.Sprintf("%s %s", errMsg, err.Error())
		}
		return spec.ReturnFail(spec.Code[spec.ExecCommandError], errMsg)
	}
	result := string(output)
	return spec.ReturnSuccess(result)
}

func isBladeCommand(script string) bool {
	return script == path.Join(util.GetProgramPath(), "blade") ||
		script == "blade"
}
