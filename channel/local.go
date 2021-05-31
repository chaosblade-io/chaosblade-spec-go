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

	"github.com/shirou/gopsutil/process"
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
			logrus.WithField("pid", p.Pid).WithError(err).Debugln("get process name error")
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
		logrus.WithFields(logrus.Fields{
			"name":        name,
			"cmdline":     cmdline,
			"processName": processName,
		}).Debugln("process info")
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
	currPid := os.Getpid()
	excludeProcesses := getExcludeProcesses(ctx)
	pids := make([]string, 0)
	for _, p := range processes {
		cmdline, err := p.Cmdline()
		if err != nil {
			logrus.WithField("pid", p.Pid).WithError(err).Debugln("get command line err")
			continue
		}
		if !strings.Contains(cmdline, processName) {
			continue
		}
		logrus.WithFields(logrus.Fields{
			"cmdline":                   cmdline,
			"processName":               processName,
			"otherConditionProcessName": otherConditionProcessName,
			"excludeProcesses":          excludeProcesses,
		}).Debugln("process info")
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

// check command is available or not
// now, all commands are: ["rm", "dd" ,"touch", "mkdir",  "echo", "kill", ,"mv","mount", "umount","tc", "head"
//"grep", "cat", "iptables", "sed", "awk", "tar"]
func (l *LocalChannel) IsAllCommandsAvailable(commandNames []string) (*spec.Response, bool) {
	if len(commandNames) == 0 {
		return nil, true
	}

	for _, commandName := range commandNames {
		if l.IsCommandAvailable(commandName) {
			continue
		}
		switch commandName {
		case "rm":
			return spec.ResponseFailWaitResult(spec.CommandRmNotFound, spec.ResponseErr[spec.CommandRmNotFound].Err,
				spec.ResponseErr[spec.CommandRmNotFound].ErrInfo), false
		case "dd":
			return spec.ResponseFailWaitResult(spec.CommandDdNotFound, spec.ResponseErr[spec.CommandDdNotFound].Err,
				spec.ResponseErr[spec.CommandDdNotFound].ErrInfo), false
		case "touch":
			return spec.ResponseFailWaitResult(spec.CommandTouchNotFound, spec.ResponseErr[spec.CommandTouchNotFound].Err,
				spec.ResponseErr[spec.CommandTouchNotFound].ErrInfo), false
		case "mkdir":
			return spec.ResponseFailWaitResult(spec.CommandMkdirNotFound, spec.ResponseErr[spec.CommandMkdirNotFound].Err,
				spec.ResponseErr[spec.CommandMkdirNotFound].ErrInfo), false
		case "echo":
			return spec.ResponseFailWaitResult(spec.CommandEchoNotFound, spec.ResponseErr[spec.CommandEchoNotFound].Err,
				spec.ResponseErr[spec.CommandEchoNotFound].ErrInfo), false
		case "kill":
			return spec.ResponseFailWaitResult(spec.CommandKillNotFound, spec.ResponseErr[spec.CommandKillNotFound].Err,
				spec.ResponseErr[spec.CommandKillNotFound].ErrInfo), false
		case "mv":
			return spec.ResponseFailWaitResult(spec.CommandMvNotFound, spec.ResponseErr[spec.CommandMvNotFound].Err,
				spec.ResponseErr[spec.CommandMvNotFound].ErrInfo), false
		case "mount":
			return spec.ResponseFailWaitResult(spec.CommandMountNotFound, spec.ResponseErr[spec.CommandMountNotFound].Err,
				spec.ResponseErr[spec.CommandMountNotFound].ErrInfo), false
		case "umount":
			return spec.ResponseFailWaitResult(spec.CommandUmountNotFound, spec.ResponseErr[spec.CommandUmountNotFound].Err,
				spec.ResponseErr[spec.CommandUmountNotFound].ErrInfo), false
		case "tc":
			return spec.ResponseFailWaitResult(spec.CommandTcNotFound, spec.ResponseErr[spec.CommandTcNotFound].Err,
				spec.ResponseErr[spec.CommandTcNotFound].ErrInfo), false
		case "head":
			return spec.ResponseFailWaitResult(spec.CommandHeadNotFound, spec.ResponseErr[spec.CommandHeadNotFound].Err,
				spec.ResponseErr[spec.CommandHeadNotFound].ErrInfo), false
		case "grep":
			return spec.ResponseFailWaitResult(spec.CommandGrepNotFound, spec.ResponseErr[spec.CommandGrepNotFound].Err,
				spec.ResponseErr[spec.CommandGrepNotFound].ErrInfo), false
		case "cat":
			return spec.ResponseFailWaitResult(spec.CommandCatNotFound, spec.ResponseErr[spec.CommandCatNotFound].Err,
				spec.ResponseErr[spec.CommandCatNotFound].ErrInfo), false
		case "iptables":
			return spec.ResponseFailWaitResult(spec.CommandIptablesNotFound, spec.ResponseErr[spec.CommandIptablesNotFound].Err,
				spec.ResponseErr[spec.CommandIptablesNotFound].ErrInfo), false
		case "sed":
			return spec.ResponseFailWaitResult(spec.CommandSedNotFound, spec.ResponseErr[spec.CommandSedNotFound].Err,
				spec.ResponseErr[spec.CommandSedNotFound].ErrInfo), false
		case "awk":
			return spec.ResponseFailWaitResult(spec.CommandAwkNotFound, spec.ResponseErr[spec.CommandAwkNotFound].Err,
				spec.ResponseErr[spec.CommandAwkNotFound].ErrInfo), false
		case "tar":
			return spec.ResponseFailWaitResult(spec.CommandAwkNotFound, spec.ResponseErr[spec.CommandAwkNotFound].Err,
				spec.ResponseErr[spec.CommandAwkNotFound].ErrInfo), false
		}
	}
	return nil, true
}

func (l *LocalChannel) IsCommandAvailable(commandName string) bool {
	response := l.Run(context.TODO(), "command", fmt.Sprintf("-v %s", commandName))
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

	pids := []string{}

	//on centos7, ss outupt pid with 'pid='
	//$ss -lpn 'sport = :80'
	//Netid State      Recv-Q Send-Q   Local Address:Port   Peer Address:Port
	//tcp   LISTEN     0      128       *:80                 *:* users:(("tengine",pid=237768,fd=6),("tengine",pid=237767,fd=6))

	//on centos6, ss output pid without 'pid='
	//$ss -lpn 'sport = :80'
	//Netid State      Recv-Q Send-Q   Local Address:Port   Peer Address:Port
	//tcp   LISTEN     0      128       *:80                 *:* users:(("tengine",237768,fd=6),("tengine",237767,fd=6))
	response := l.Run(context.TODO(), "ss", fmt.Sprintf("-pln sport = :%s", localPort))
	if !response.Success {
		return pids, fmt.Errorf(response.Err)
	}
	if util.IsNil(response.Result) {
		return pids, nil
	}
	result := response.Result.(string)
	ssMsg := strings.TrimSpace(result)
	if ssMsg == "" {
		return pids, nil
	}
	sockets := strings.Split(ssMsg, "\n")
	logrus.Infof("sockets for %s, %v", localPort, sockets)
	for idx, s := range sockets {
		if idx == 0 {
			continue
		}
		fields := strings.Fields(s)
		// centos7: users:(("tengine",pid=237768,fd=6),("tengine",pid=237767,fd=6))
		// centos6: users:(("tengine",237768,fd=6),("tengine",237767,fd=6))
		lastField := fields[len(fields)-1]
		logrus.Infof("GetPidsByLocalPort: lastField: %v", lastField)
		pidExp := regexp.MustCompile(`pid=(\d+)|,(\d+),`)
		// extract all the pids that conforms to pidExp
		matchedPidArrays := pidExp.FindAllStringSubmatch(lastField, -1)
		if matchedPidArrays == nil || len(matchedPidArrays) == 0 {
			return pids, nil
		}

		for _, matchedPidArray := range matchedPidArrays {

			var pid string

			// centos7: matchedPidArray is [pid=29863 29863 ], matchedPidArray[len(matchedPidArray)-1] is whitespace

			pid = strings.TrimSpace(matchedPidArray[len(matchedPidArray)-1])

			if pid != "" {
				pids = append(pids, pid)
				continue
			}

			// centos6: matchedPidArray is [,237768,  237768] matchedPidArray[len(matchedPidArray)-1] is pid
			pid = strings.TrimSpace(matchedPidArray[len(matchedPidArray)-2])
			if pid != "" {
				pids = append(pids, pid)
				continue
			}

		}
	}
	logrus.Infof("GetPidsByLocalPort: pids: %v", pids)
	return pids, nil
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
	if resp := isBinBladeCommand(script); resp != nil {
		return resp
	}

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script+" "+args)
	output, err := cmd.CombinedOutput()
	outMsg := string(output)

	if err == nil {
		return spec.ReturnSuccess(outMsg)
	}

	if strings.Contains(outMsg, "RTNETLINK answers: File exists") {
		return spec.ResponseFail(spec.CommandNetworkExist, spec.ResponseErr[spec.CommandNetworkExist].ErrInfo)
	}
	if !isBladeCmd {
		outMsg += " " + err.Error()
	}
	return spec.ResponseFail(spec.OsCmdExecFailed, fmt.Sprintf(spec.ResponseErr[spec.OsCmdExecFailed].ErrInfo, script+" "+args, outMsg))
}
func isBinBladeCommand(script string) *spec.Response {
	if ok := strings.Contains(script, "/chaosblade"); !ok {
		return nil
	}

	if ok := util.IsExist(script); !ok {
		return spec.ResponseFail(spec.ChaosbladeFileNotFound, fmt.Sprintf(spec.ResponseErr[spec.ChaosbladeFileNotFound].Err, script))
	}
	return nil
}

func isBladeCommand(script string) bool {
	return script == path.Join(util.GetProgramPath(), "blade") ||
		script == "blade"
}
