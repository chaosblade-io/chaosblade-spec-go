//go:build windows
// +build windows

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
	"github.com/chaosblade-io/chaosblade-spec-go/log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"github.com/shirou/gopsutil/process"
)

type LocalChannel struct {
}

// NewLocalChannel returns a local channel for invoking the host command
func NewLocalChannel() spec.Channel {
	return &LocalChannel{}
}

func (l *LocalChannel) Name() string {
	return "local"
}

func (l *LocalChannel) Run(ctx context.Context, script, args string) *spec.Response {
	return execScript(ctx, script, args)
}

func (l *LocalChannel) GetScriptPath() string {
	return util.GetProgramPath()
}

func (l *LocalChannel) GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error) {
	processName = strings.TrimSpace(processName)
	if processName == "" {
		return []string{}, fmt.Errorf("processName is blank")
	}
	processes, err := process.Processes()
	if err != nil {
		return []string{}, err
	}
	currPid := os.Getpid()
	excludeProcesses := getExcludeProcesses(ctx)
	pids := make([]string, 0)
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			log.Debugf(ctx, "get process name error, pid: %d, err: %v", p.Pid, err)
			continue
		}
		if processName != name {
			continue
		}
		if int32(os.Getpid()) == p.Pid {
			continue
		}
		cmdline, _ := p.Cmdline()
		containsExcludeProcess := false
		log.Debugf(ctx, "process info, name: %s, cmdline: %s, processName: %s", name, cmdline, processName)
		for _, ep := range excludeProcesses {
			if strings.Contains(cmdline, strings.TrimSpace(ep)) {
				containsExcludeProcess = true
				break
			}
		}
		if containsExcludeProcess {
			continue
		}
		if p.Pid == int32(currPid) {
			continue
		}
		pids = append(pids, fmt.Sprintf("%d", p.Pid))
	}
	return pids, nil
}

func (l *LocalChannel) GetPidsByProcessName(processName string, ctx context.Context) ([]string, error) {
	processName = strings.TrimSpace(processName)
	if processName == "" {
		return []string{}, fmt.Errorf("process keyword is blank")
	}
	processes, err := process.Processes()
	if err != nil {
		return []string{}, err
	}
	otherConditionProcessValue := ctx.Value(ProcessKey)
	otherConditionProcessName := ""
	if otherConditionProcessValue != nil {
		otherConditionProcessName = otherConditionProcessValue.(string)
	}
	processCommandValue := ctx.Value(ProcessCommandKey)
	processCommandName := ""
	if processCommandValue != nil {
		processCommandName = processCommandValue.(string)
	}
	currPid := os.Getpid()
	excludeProcesses := getExcludeProcesses(ctx)
	pids := make([]string, 0)
	for _, p := range processes {
		if processCommandName != "" {
			name, err := p.Name()
			if err != nil {
				log.Debugf(ctx, "get process command error, processCommand: %s, err: %v, ", processCommandName, err)
				continue
			}
			if !strings.Contains(name, processCommandName) {
				continue
			}
		}
		cmdline, err := p.Cmdline()
		if err != nil {
			log.Debugf(ctx, "get command line error, pid: %d, err: %v", p.Pid, err)
			continue
		}
		if !strings.Contains(cmdline, processName) {
			continue
		}
		log.Debugf(ctx, "process info, cmdline: %s, processName: %s, processCommand: %s, otherConditionProcessName: %s, excludeProcesses: %s",
			cmdline, processName, processCommandName, otherConditionProcessName, excludeProcesses)

		if otherConditionProcessName != "" && !strings.Contains(cmdline, otherConditionProcessName) {
			continue
		}
		containsExcludeProcess := false
		for _, ep := range excludeProcesses {
			if strings.Contains(cmdline, ep) {
				containsExcludeProcess = true
				break
			}
		}
		if containsExcludeProcess {
			continue
		}
		if p.Pid == int32(currPid) {
			continue
		}
		pids = append(pids, fmt.Sprintf("%d", p.Pid))
	}
	return pids, nil
}

func getExcludeProcesses(ctx context.Context) []string {
	excludeProcessValue := ctx.Value(ExcludeProcessKey)
	excludeProcesses := make([]string, 0)
	if excludeProcessValue != nil {
		excludeProcessesString := excludeProcessValue.(string)
		processNames := strings.Split(excludeProcessesString, ",")
		for _, name := range processNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			excludeProcesses = append(excludeProcesses, name)
		}
	}
	excludeProcesses = append(excludeProcesses, "chaos_killprocess", "chaos_stopprocess")
	return excludeProcesses
}

func (l *LocalChannel) GetPsArgs(ctx context.Context) string {
	var psArgs = "-eo user,pid,ppid,args"
	if l.IsAlpinePlatform(ctx) {
		psArgs = "-o user,pid,ppid,args"
	}
	return psArgs
}

func (l *LocalChannel) IsAlpinePlatform(ctx context.Context) bool {
	var osVer = ""
	if util.IsExist("/etc/os-release") {
		response := l.Run(ctx, "awk", "-F '=' '{if ($1 == \"ID\") {print $2;exit 0}}' /etc/os-release")
		if response.Success {
			osVer = response.Result.(string)
		}
	}
	return strings.TrimSpace(osVer) == "alpine"
}

// check command is available or not
// now, all commands are: ["rm", "dd" ,"touch", "mkdir",  "echo", "kill", ,"mv","mount", "umount","tc", "head"
//"grep", "cat", "iptables", "sed", "awk", "tar"]
func (l *LocalChannel) IsAllCommandsAvailable(ctx context.Context, commandNames []string) (*spec.Response, bool) {
	return IsAllCommandsAvailable(ctx, l, commandNames)
}

func (l *LocalChannel) IsCommandAvailable(ctx context.Context, commandName string) bool {
	response := l.Run(ctx, "command", fmt.Sprintf("-v %s", commandName))
	return response.Success
}

func (l *LocalChannel) ProcessExists(pid string) (bool, error) {
	p, err := strconv.Atoi(pid)
	if err != nil {
		return false, err
	}
	return process.PidExists(int32(p))
}

func (l *LocalChannel) GetPidUser(pid string) (string, error) {
	p, err := strconv.Atoi(pid)
	if err != nil {
		return "", err
	}
	process, err := process.NewProcess(int32(p))
	if err != nil {
		return "", err
	}
	return process.Username()
}

func (l *LocalChannel) GetPidsByLocalPorts(ctx context.Context, localPorts []string) ([]string, error) {
	if localPorts == nil || len(localPorts) == 0 {
		return nil, fmt.Errorf("the local port parameter is empty")
	}
	var result = make([]string, 0)
	for _, port := range localPorts {
		pids, err := l.GetPidsByLocalPort(ctx, port)
		if err != nil {
			return nil, fmt.Errorf("failed to get pid by %s, %v", port, err)
		}
		log.Infof(ctx, "get pids by %s port returns %v", port, pids)
		if pids != nil && len(pids) > 0 {
			result = append(result, pids...)
		}
	}
	return result, nil
}

func (l *LocalChannel) GetPidsByLocalPort(ctx context.Context, localPort string) ([]string, error) {
	return GetPidsByLocalPort(ctx, l, localPort)
}

// execScript invokes exec.CommandContext
func execScript(ctx context.Context, script, args string) *spec.Response {
	isBladeCommand := isBladeCommand(script)
	if isBladeCommand && !util.IsExist(script) {
		// TODO nohup invoking
		return spec.ResponseFailWithFlags(spec.ChaosbladeFileNotFound, script)
	}
	newCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if ctx == context.Background() {
		ctx = newCtx
	}
	log.Debugf(ctx, "Command: %s %s", script, args)
	cmd := exec.CommandContext(ctx, "cmd", "/C", script+` `+args)
	output, err := cmd.CombinedOutput()
	outMsg := string(output)
	log.Debugf(ctx, "Command Result, output: %v, err: %v", outMsg, err)
	if strings.TrimSpace(outMsg) != "" {
		resp := spec.Decode(outMsg, nil)
		if resp.Code != spec.ResultUnmarshalFailed.Code {
			return resp
		}
	}
	if err == nil {
		return spec.ReturnSuccess(outMsg)
	}
	outMsg += " " + err.Error()
	return spec.ResponseFailWithFlags(spec.OsCmdExecFailed, cmd, outMsg)
}

func isBladeCommand(script string) bool {
	return strings.HasSuffix(script, util.GetProgramPath())
}
