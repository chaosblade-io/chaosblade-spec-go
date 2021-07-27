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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type CodeType struct {
	Code int32
	Msg  string
}

var (
	IgnoreCode                        = CodeType{100, "ignore code"}
	OK                                = CodeType{200, "success"}
	ReturnOKDirectly                  = CodeType{201, "return ok directly"}
	Forbidden                         = CodeType{43000, "Forbidden: must be root"}
	ActionNotSupport                  = CodeType{44000, "`%s`: action not supported"}
	ParameterLess                     = CodeType{45000, "less parameter: `%s`"}
	ParameterIllegal                  = CodeType{46000, "illegal `%s` parameter value: `%s`. %v"}
	ParameterInvalid                  = CodeType{47000, "invalid `%s` parameter value: `%s`. %v"}
	ParameterInvalidProName           = CodeType{47001, "invalid parameter `%s`, `%s` process not found"}
	ParameterInvalidProIdNotByName    = CodeType{47002, "invalid parameter `process|pid`, the process ids got by %s does not contain the pid %s value"}
	ParameterInvalidCplusPort         = CodeType{47003, "invalid parameter port, `%s` port not found, please execute prepare command firstly"}
	ParameterInvalidDbQuery           = CodeType{47004, "invalid parameter `%s`, db record not found"}
	ParameterInvalidCplusTarget       = CodeType{47005, "invalid parameter target, `%s` target not support"}
	ParameterInvalidBladePathError    = CodeType{47006, "invalid parameter `%s`, deploy chaosblade to `%s` failed, err: %v"}
	ParameterInvalidNSNotOne          = CodeType{47007, "invalid parameter `%s`, only one value can be specified"}
	ParameterInvalidK8sPodQuery       = CodeType{47008, "invalid parameter `%s`, can not find pods"}
	ParameterInvalidK8sNodeQuery      = CodeType{47009, "invalid parameter `%s`, can not find node"}
	ParameterInvalidDockContainerId   = CodeType{47010, "invalid parameter `%s`, can not find container by id"}
	ParameterInvalidDockContainerName = CodeType{47011, "invalid parameter `%s`, can not find container by name"}
	ParameterRequestFailed            = CodeType{48000, "get request parameter failed"}
	CommandIllegal                    = CodeType{49000, "illegal command, err: %v"}
	CommandNetworkExist               = CodeType{49001, "network tc exec failed! RTNETLINK answers: File exists"}
	ChaosbladeFileNotFound            = CodeType{51000, "`%s`: chaosblade file not found"}
	CommandTasksetNotFound            = CodeType{52000, "`taskset`: command not found"}
	CommandMountNotFound              = CodeType{52001, "`mount`: command not found"}
	CommandUmountNotFound             = CodeType{52002, "`umount`: command not found"}
	CommandTcNotFound                 = CodeType{52003, "`tc`: command not found"}
	CommandIptablesNotFound           = CodeType{52004, "`iptables`: command not found"}
	CommandSedNotFound                = CodeType{52005, "`sed`: command not found"}
	CommandCatNotFound                = CodeType{52006, "`cat`: command not found"}
	CommandSsNotFound                 = CodeType{52007, "`ss`: command not found"}
	CommandDdNotFound                 = CodeType{52008, "`dd`: command not found"}
	CommandRmNotFound                 = CodeType{52009, "`rm`: command not found"}
	CommandTouchNotFound              = CodeType{52010, "`touch`: command not found"}
	CommandMkdirNotFound              = CodeType{52011, "`mkdir`: command not found"}
	CommandEchoNotFound               = CodeType{52012, "`echo`: command not found"}
	CommandKillNotFound               = CodeType{52013, "`kill`: command not found"}
	CommandMvNotFound                 = CodeType{52014, "`mv`: command not found"}
	CommandHeadNotFound               = CodeType{52015, "`head`: command not found"}
	CommandGrepNotFound               = CodeType{52016, "`grep`: command not found"}
	CommandAwkNotFound                = CodeType{52017, "`awk`: command not found"}
	CommandTarNotFound                = CodeType{52018, "`tar`: command not found"}
	CommandSystemctlNotFound          = CodeType{52019, "`systemctl`: command not found"}
	CommandNohupNotFound              = CodeType{52020, "`nohup`: command not found"}
	ChaosbladeServerStarted           = CodeType{53000, "the chaosblade has been started. If you want to stop it, you can execute blade server stop command"}
	UnexpectedStatus                  = CodeType{54000, "unexpected status, expected status: `%s`, but the real status: `%s`, please wait!"}
	DockerExecNotFound                = CodeType{55000, "`%s`: the docker exec not found"}
	DockerImagePullFailed             = CodeType{55001, "pull image failed, err: %v"}
	HandlerExecNotFound               = CodeType{56000, "`%s`: the handler exec not found"}
	CplusActionNotSupport             = CodeType{56001, "`%s`: cplus action not support"}
	ContainerInContextNotFound        = CodeType{56002, "cannot find container, please confirm if the container exists"}
	PodNotReady                       = CodeType{56003, "`%s` pod is not ready"}
	ResultUnmarshalFailed             = CodeType{60000, "`%s`: exec result unmarshal failed, err: %v"}
	ResultMarshalFailed               = CodeType{60001, "`%v`: exec result marshal failed, err: %v"}
	GenerateUidFailed                 = CodeType{60002, "generate experiment uid failed, err: %v"}
	ChaosbladeServiceStoped           = CodeType{61000, "chaosblade service has been stopped"}
	ProcessIdByNameFailed             = CodeType{63010, "`%s`: get process id by name failed, err: %v"}
	ProcessJudgeExistFailed           = CodeType{63011, "`%s`: judge the process exist or not, failed, err: %v"}
	ProcessNotExist                   = CodeType{63012, "`%s`: the process not exist"}
	ProcessGetUsernameFailed          = CodeType{63014, "`%s`: get username failed by the process id, err: %v"}
	ChannelNil                        = CodeType{63020, "chanel is nil"}
	SandboxGetPortFailed              = CodeType{63030, "get sandbox port failed, err: %v"}
	SandboxCreateTokenFailed          = CodeType{63031, "create sandbox token failed, err: %v"}
	FileCantGetLogFile                = CodeType{63040, "can not get log file"}
	FileNotExist                      = CodeType{63041, "`%s`: not exist"}
	FileCantReadOrOpen                = CodeType{63042, "`%s`: can not read or open"}
	BackfileExists                    = CodeType{63050, "`%s`: backup file exists, may be annother experiment is running"}
	DbQueryFailed                     = CodeType{63060, "`%s`: db query failed, err: %v"}
	K8sExecFailed                     = CodeType{63061, "`%s`: k8s exec failed, err: %v"}
	DockerExecFailed                  = CodeType{63062, "`%s`: docker exec failed, err: %v"}
	OsCmdExecFailed                   = CodeType{63063, "`%s`: cmd exec failed, err: %v"}
	HttpExecFailed                    = CodeType{63064, "`%s`: http cmd failed, err: %v"}
	GetIdentifierFailed               = CodeType{63065, "get experiment identifier failed, err: %v"}
	OsExecutorNotFound                = CodeType{63070, "`%s`: os executor not found"}
	ChaosfsClientFailed               = CodeType{64000, "init chaosfs client failed in pod %v, err: %v"}
	ChaosfsInjectFailed               = CodeType{64001, "inject io exception in pod %s failed, request %v, err: %v"}
	ChaosfsRecoverFailed              = CodeType{64002, "recover io exception failed in pod  %v, err: %v"}
	SshExecFailed                     = CodeType{65000, "ssh exec failed, result: %v, err %v"}
	SshExecNothing                    = CodeType{65001, "cannot get result from remote host, please execute recovery and try again"}
	SystemdNotFound                   = CodeType{66001, "`%s`: systemd not found, err: %v"}
	DatabaseError                     = CodeType{67001, "`%s`: failed to execute, err: %v"}
	DataNotFound                      = CodeType{67002, "`%s` record not found, if it's k8s experiment, please add --target k8s flag to retry"}
)

func (c CodeType) Sprintf(values ...interface{}) string {
	return fmt.Sprintf(c.Msg, values...)
}

type Response struct {
	Code    int32       `json:"code"`
	Success bool        `json:"success"`
	Err     string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func (response *Response) Error() string {
	return response.Print()
}

func (response *Response) Print() string {
	bytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf("marshall response err, %s; code: %d", err.Error(), response.Code)
	}
	return string(bytes)
}

func Return(codeType CodeType, success bool) *Response {
	return &Response{Code: codeType.Code, Success: success, Err: codeType.Msg}
}

func ReturnFail(codeType CodeType, err string) *Response {
	return &Response{Code: codeType.Code, Success: false, Err: err}
}

func ReturnSuccess(result interface{}) *Response {
	return &Response{Code: OK.Code, Success: true, Result: result}
}

func ReturnResultIgnoreCode(result interface{}) *Response {
	return &Response{Code: IgnoreCode.Code, Result: result}
}

func ResponseFail(status int32, err string, result interface{}) *Response {
	return &Response{Code: status, Success: false, Err: err, Result: result}
}

func ResponseFailWithFlags(codeType CodeType, flags ...interface{}) *Response {
	if flags == nil {
		return &Response{Code: codeType.Code, Success: false, Err: codeType.Msg}
	}
	return &Response{Code: codeType.Code, Success: false, Err: fmt.Sprintf(codeType.Msg, flags...)}
}

func ResponseFailWithResult(codeType CodeType, result interface{}, flags ...interface{}) *Response {
	return &Response{Code: codeType.Code, Success: false, Result: result, Err: fmt.Sprintf(codeType.Msg, flags...)}
}

func Success() *Response {
	return ReturnSuccess(nil)
}

//ToString
func (response *Response) ToString() string {
	bytes, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintln(string(bytes))
}

// Decode return the response that wraps the content
func Decode(content string, defaultValue *Response) *Response {
	var resp Response
	content = strings.TrimSpace(content)
	err := json.Unmarshal([]byte(content), &resp)
	if err != nil {
		if defaultValue == nil {
			defaultValue = ResponseFailWithFlags(ResultUnmarshalFailed, content, err.Error())
		}
		logrus.Debugf("decode %s err, return default value, %s", content, defaultValue.Print())
		return defaultValue
	}
	return &resp
}
