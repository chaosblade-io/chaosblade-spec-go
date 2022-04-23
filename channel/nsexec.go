package channel

import (
	"context"
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/log"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	NSTargetFlagName = "ns_target"
	NSPidFlagName    = "ns_pid"
	NSMntFlagName    = "ns_mnt"
	NSNetFlagName    = "ns_net"
)

type NSExecChannel struct {
	LocalChannel
}

func NewNSExecChannel() spec.Channel {
	return &NSExecChannel{}
}

func (l *NSExecChannel) Name() string {
	return "nsexec"
}

func (l *NSExecChannel) Run(ctx context.Context, script, args string) *spec.Response {
	pid := ctx.Value(NSTargetFlagName)
	if pid == nil {
		return spec.ResponseFailWithFlags(spec.CommandIllegal, script)
	}

	ns_script := fmt.Sprintf("-t %s", pid)

	if ctx.Value(NSPidFlagName) == spec.True {
		ns_script = fmt.Sprintf("%s -p", ns_script)
	}

	if ctx.Value(NSMntFlagName) == spec.True {
		ns_script = fmt.Sprintf("%s -m", ns_script)
	}

	if ctx.Value(NSNetFlagName) == spec.True {
		ns_script = fmt.Sprintf("%s -n", ns_script)
	}

	isBladeCommand := isBladeCommand(script)
	if isBladeCommand && !util.IsExist(script) {
		// TODO nohup invoking
		return spec.ResponseFailWithFlags(spec.ChaosbladeFileNotFound, script)
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	if args != "" {
		args = script + " " + args
	} else {
		args = script
	}

	ns_script = fmt.Sprintf("%s -- /bin/sh -c", ns_script)

	programPath := util.GetProgramPath()
	if path.Base(programPath) != spec.BinPath {
		programPath = path.Join(programPath, spec.BinPath)
	}
	bin := path.Join(programPath, spec.NSExecBin)
	log.Debugf(ctx,`Command: %s %s "%s"`, bin, ns_script, args)

	split := strings.Split(ns_script, " ")

	cmd := exec.CommandContext(timeoutCtx, bin, append(split, args)...)
	output, err := cmd.CombinedOutput()
	outMsg := string(output)
	log.Debugf(ctx, "Command Result, output: %v, err: %v", outMsg, err)
	// TODO shell-init错误
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

func (l *NSExecChannel) GetPidsByProcessCmdName(processName string, ctx context.Context) ([]string, error) {
	excludeProcesses := ctx.Value(ExcludeProcessKey)
	excludeGrepInfo := ""
	if excludeProcesses != nil {
		excludeProcessesString := excludeProcesses.(string)
		excludeProcessArrays := strings.Split(excludeProcessesString, ",")
		for _, excludeProcess := range excludeProcessArrays {
			if excludeProcess != "" {
				excludeGrepInfo += fmt.Sprintf(`| grep -v -w %s`, excludeProcess)
			}
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

func (l *NSExecChannel) GetPidsByProcessName(processName string, ctx context.Context) ([]string, error) {
	psArgs := l.GetPsArgs(ctx)
	otherProcess := ctx.Value(ProcessKey)
	otherGrepInfo := ""
	if otherProcess != nil {
		processString := otherProcess.(string)
		if processString != "" {
			otherGrepInfo = fmt.Sprintf(`| grep "%s"`, processString)
		}
	}
	excludeProcesses := ctx.Value(ExcludeProcessKey)
	excludeGrepInfo := ""
	if excludeProcesses != nil {
		excludeProcessesString := excludeProcesses.(string)
		excludeProcessArrays := strings.Split(excludeProcessesString, ",")
		for _, excludeProcess := range excludeProcessArrays {
			if excludeProcess != "" {
				excludeGrepInfo += fmt.Sprintf(`| grep -v -w %s`, excludeProcess)
			}
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

func (l *NSExecChannel) IsAllCommandsAvailable(ctx context.Context, commandNames []string) (*spec.Response, bool) {
	return IsAllCommandsAvailable(ctx, l, commandNames)
}

func (l *NSExecChannel) IsCommandAvailable(ctx context.Context, commandName string) bool {
	response := l.Run(ctx, "command", fmt.Sprintf("-v %s", commandName))
	if response.Success {
		if response.Result != nil && strings.Contains(response.Result.(string), commandName) {
			return true
		}
	}
	return false
}

func (l *NSExecChannel) GetPsArgs(ctx context.Context) string {
	var psArgs = "-eo user,pid,ppid,args"
	if l.IsAlpinePlatform(ctx) {
		psArgs = "-o user,pid,ppid,args"
	}
	return psArgs
}

func (l *NSExecChannel) IsAlpinePlatform(ctx context.Context) bool {
	var osVer = ""
	if util.IsExist("/etc/os-release") {
		response := l.Run(ctx, "awk", "-F '=' '{if ($1 == \"ID\") {print $2;exit 0}}' /etc/os-release")
		if response.Success {
			osVer = response.Result.(string)
		}
	}
	return strings.TrimSpace(osVer) == "alpine"
}

func (l *NSExecChannel) GetPidsByLocalPort(ctx context.Context, localPort string) ([]string, error) {
	return GetPidsByLocalPort(ctx, l, localPort)
}
