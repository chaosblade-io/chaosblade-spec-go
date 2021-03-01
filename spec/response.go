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

	"github.com/sirupsen/logrus"
)

const (
	IgnoreCode       = "IgnoreCode"
	OK               = "OK"
	ReturnOKDirectly = "ReturnOKDirectly"
	InvalidTimestamp = "InvalidTimestamp"
	//Forbidden               = "Forbidden"
	HandlerNotFound         = "HandlerNotFound"
	TokenNotFound           = "TokenNotFound"
	DataNotFound            = "DataNotFound"
	GetProcessError         = "GetProcessError"
	ServerError             = "ServerError"
	HandlerClosed           = "HandlerClosed"
	Timeout                 = "Timeout"
	Uninitialized           = "Uninitialized"
	EncodeError             = "EncodeError"
	DecodeError             = "DecodeError"
	FileNotFound            = "FileNotFound"
	DownloadError           = "DownloadError"
	DeployError             = "DeployError"
	ServiceSwitchError      = "ServiceSwitchError"
	DiskNotFound            = "DiskNotFound"
	DatabaseError           = "DatabaseError"
	EnvironmentError        = "EnvironmentError"
	NoWritePermission       = "NoWritePermission"
	RemoveRecordError       = "RemoveRecordError"
	ParameterEmpty          = "ParameterEmpty"
	ParameterTypeError      = "ParameterTypeError"
	IllegalParameters       = "IllegalParameters"
	IllegalCommand          = "IllegalCommand"
	ExecCommandError        = "ExecCommandError"
	DuplicateError          = "DuplicateError"
	FaultInjectCmdError     = "FaultInjectCmdError"
	FaultInjectExecuteError = "FaultInjectExecuteError"
	FaultInjectNotSupport   = "FaultInjectNotSupport"
	JavaAgentCmdError       = "JavaAgentCmdError"
	K8sInvokeError          = "K8sInvokeError"
	DockerInvokeError       = "DockerInvokeError"
	DestroyNotSupported     = "DestroyNotSupported"
	PreHandleError          = "PreHandleError"
	SandboxInvokeError      = "SandboxInvokeError"
	CommandNotFound         = "CommandNotFound"
	StatusError             = "StatusError"
	UnexpectedCommandError  = "UnexpectedCommandError"
	CplusProxyCmdError      = "CplusProxyCmdError"
)

type CodeType struct {
	Code int32
	Msg  string
}

var Code = map[string]CodeType{
	IgnoreCode:       {100, "ignore code"},
	OK:               {200, "success"},
	ReturnOKDirectly: {201, "return ok directly"},
	InvalidTimestamp: {401, "invalid timestamp"},
	//Forbidden:               {403, "forbidden"},
	HandlerNotFound:         {404, "request handler not found"},
	TokenNotFound:           {405, "access token not found"},
	DataNotFound:            {406, "data not found"},
	DestroyNotSupported:     {407, "destroy not supported"},
	GetProcessError:         {408, "get process error"},
	ServerError:             {500, "server error"},
	HandlerClosed:           {501, "handler closed"},
	PreHandleError:          {502, "pre handle error"},
	CommandNotFound:         {503, "command not found"},
	StatusError:             {504, "status error"},
	Timeout:                 {510, "timeout"},
	Uninitialized:           {511, "uninitialized"},
	EncodeError:             {512, "encode error"},
	DecodeError:             {513, "decode error"},
	FileNotFound:            {514, "file not found"},
	DownloadError:           {515, "download file error"},
	DeployError:             {516, "deploy file error"},
	ServiceSwitchError:      {517, "service switch error"},
	DiskNotFound:            {518, "disk not found"},
	DatabaseError:           {520, "execute db error"},
	EnvironmentError:        {521, "environment error"},
	NoWritePermission:       {522, "no write permission"},
	RemoveRecordError:       {530, "remove record or resource err"},
	ParameterEmpty:          {600, "parameter is empty"},
	ParameterTypeError:      {601, "parameter type error"},
	IllegalParameters:       {602, "illegal parameters"},
	IllegalCommand:          {603, "illegal command"},
	ExecCommandError:        {604, "exec command error"},
	DuplicateError:          {605, "duplicate error"},
	FaultInjectCmdError:     {701, "cannot handle the faultInject cmd"},
	FaultInjectExecuteError: {702, "execute faultInject error"},
	FaultInjectNotSupport:   {703, "the inject type not support"},
	JavaAgentCmdError:       {704, "cannot handle the javaagent cmd"},
	K8sInvokeError:          {800, "invoke k8s server api error"},
	DockerInvokeError:       {801, "invoke docker command error"},
	SandboxInvokeError:      {802, "invoke sandbox error"},
	CplusProxyCmdError:      {803, "invoke cplus proxy error"},
	UnexpectedCommandError:  {901, "unexpected command error"},
}

type ResultType struct {
	Err     string
	ErrInfo string
}

const (
	Success = 200
	// 2. failed
	// 2.1 client error
	//Uninitialized          = 41000
	Forbidden                       = 43000
	ActionNotSupport                = 44000
	ParameterLess                   = 45000
	ParameterIllegal                = 46000
	ParameterInvalid                = 47000
	ParameterInvalidCplusPort       = 47001
	ParameterInvalidProName         = 47002
	ParameterInvalidProIdNotByName  = 47003
	ParameterInvalidDbQuery         = 47004
	ParameterInvalidCplusTarget     = 47005
	ParameterInvalidBladePathError  = 47006
	ParameterInvalidNSNotOne        = 47007
	ParameterInvalidK8sPodQuery     = 47008
	ParameterInvalidK8sNodeQuery    = 47009
	ParameterInvalidDockContainerId = 47010
	ParameterRequestFailed          = 48000
	CommandLess                     = 49000
	CommandNetworkExist             = 49001

	// 2.2 server error, but the user can hold it
	ChaosbladeFileNotFound  = 51000
	CommandTasksetNotFound  = 52000
	CommandMountNotFound    = 52001
	CommandUmountNotFound   = 52002
	CommandTcNotFound       = 52003
	CommandIptablesNotFound = 52004
	CommandSedNotFound      = 52005
	CommandCatNotFound      = 52006
	CommandSsNotFound       = 52007
	CommandDdNotFound       = 52008
	CommandRmNotFound       = 52009
	CommandTouchNotFound    = 52010
	CommandMkdirNotFound    = 52011
	CommandEchoNotFound     = 52012
	CommandKillNotFound     = 52013
	CommandMvNotFound       = 52014
	CommandHeadNotFound     = 52015
	CommandGrepNotFound     = 52016
	CommandAwkNotFound      = 52017
	CommandTarNotFound      = 52018
	ChaosbladeServerStarted = 53000
	UnexpectedStatus        = 54000
	DockerExecNotFound      = 55000
	DockerImagePullFailed   = 55001
	HandlerExecNotFound     = 56000
	CplusActionNotSupport   = 56001

	// 2.3 server error, but the user can not hold it
	ResultUnmarshalFailed    = 60000
	ResultMarshalFailed      = 60001
	ChaosbladeServiceStoped  = 61000
	ProcessIdByNameFailed    = 63010
	ProcessJudgeExistFailed  = 63011
	ProcessNotExist          = 63012
	ProcessGetUsernameFailed = 63014
	ChannelNil               = 63020
	SandboxGetPortFailed     = 63030
	SandboxCreateTokenFailed = 63031
	FileCantGetLogFile       = 63040
	FileNotExist             = 63041
	FileCantReadOrOpen       = 63042
	BackfileExists           = 63050
	DbQueryFailed            = 63060
	K8sExecFailed            = 63061
	DockerExecFailed         = 63062
	OsCmdExecFailed          = 63063
	HttpExecFailed           = 63064
	OsExecutorNotFound       = 63070
	ChaosfsClientFailed      = 64000
	ChaosfsInjectFailed      = 64001
	ChaosfsRecoverFailed     = 64002
	SshExecFailed            = 65000
	SshExecNothing           = 65001
)

var ResponseErr = map[int32]ResultType{
	//Uninitialized:       {"Uninitialized: access token not found", "Uninitialized: access token not found"},
	Forbidden:                       {"Forbidden: must be root", "Forbidden: must be root"},
	ActionNotSupport:                {"`%s`: action not supported", "`%s`: action not supported"},
	ParameterLess:                   {"less parameter: `%s`", "less parameter: `%s`"},
	ParameterIllegal:                {"illegal parameter: `%s`", "illegal parameter: `%s`"},
	ParameterInvalid:                {"invalid parameter: `%s`", "invalid parameter: `%s`"},
	ParameterInvalidProName:         {"invalid parameter `%s`, `%s` process not found", "invalid parameter `%s`, `%s` process not found"},
	ParameterInvalidProIdNotByName:  {"invalid parameter `process|pid`, the process ids got by `%s` does not contain the pid `%s` value", "invalid parameter `process|pid`, the process ids got by %s does not contain the pid %s value"},
	ParameterInvalidCplusPort:       {"invalid parameter `port`, `%s` port not found, please execute prepare command firstly", "invalid parameter port, `%s` port not found, please execute prepare command firstly"},
	ParameterInvalidDbQuery:         {"invalid parameter `%s`, db record not found", "invalid parameter `%s`, db record not found"},
	ParameterInvalidCplusTarget:     {"invalid parameter `target`, `%s` target not support", "invalid parameter target, `%s` target not support"},
	ParameterInvalidBladePathError:  {"invalid parameter `install-path`, chaosblade install into `%s` failed, err: %v", "invalid parameter `install-path`, chaosblade install into `%s` failed, err: %v"},
	ParameterInvalidNSNotOne:        {"invalid parameter `%s`, only one value can be specified", "invalid parameter `%s`, only one value can be specified"},
	ParameterInvalidK8sPodQuery:     {"invalid parameter `%s`, can not find pods", "invalid parameter `%s`, can not find pods"},
	ParameterInvalidK8sNodeQuery:    {"invalid parameter `%s`, can not find node", "invalid parameter `%s`, can not find node"},
	ParameterInvalidDockContainerId: {"invalid parameter `%s`, can not find container by id", "invalid parameter `%s`, can not find container by id"},
	ParameterRequestFailed:          {"get request parameter failed", "get request parameter failed"},
	CommandLess:                     {"less target command", "less target command"},
	CommandNetworkExist:             {"network tc exec failed! RTNETLINK answers: File exists", "network tc exec failed! RTNETLINK answers: File exists"},

	ChaosbladeFileNotFound:  {"`%s`: chaosblade file not found", "`%s`: chaosblade file not found"},
	CommandTasksetNotFound:  {"`taskset`: command not found", "`taskset`: command not found"},
	CommandMountNotFound:    {"`mount`: command not found", "`mount`: command not found"},
	CommandUmountNotFound:   {"`umount`: command not found", "`umount`: command not found"},
	CommandTcNotFound:       {"`tc`: command not found", "`tc`: command not found"},
	CommandIptablesNotFound: {"`iptables`: command not found", "`iptables`: command not found"},
	CommandSedNotFound:      {"`sed`: command not found", "`sed`: command not found"},
	CommandCatNotFound:      {"`cat`: command not found", "`cat`: command not found"},
	CommandSsNotFound:       {"`ss`: command not found", "`ss`: command not found"},
	CommandDdNotFound:       {"`dd`: command not found", "`dd`: command not found"},
	CommandRmNotFound:       {"`rm`: command not found", "`rm`: command not found"},
	CommandTouchNotFound:    {"`touch`: command not found", "`touch`: command not found"},
	CommandMkdirNotFound:    {"`mkdir`: command not found", "`mkdir`: command not found"},
	CommandEchoNotFound:     {"`echo`: command not found", "`echo`: command not found"},
	CommandKillNotFound:     {"`kill`: command not found", "`kill`: command not found"},
	CommandMvNotFound:       {"`mv`: command not found", "`mv`: command not found"},
	CommandHeadNotFound:     {"`head`: command not found", "`head`: command not found"},
	CommandGrepNotFound:     {"`grep`: command not found", "`grep`: command not found"},
	CommandAwkNotFound:      {"`awk`: command not found", "`awk`: command not found"},
	CommandTarNotFound:      {"`tar`: command not found", "`tar`: command not found"},
	ChaosbladeServerStarted: {"the chaosblade has been started", "the chaosblade has been started. If you want to stop it, you can execute blade server stop command"},
	UnexpectedStatus:        {"unexpected status, expected status: `%s`, but the real status: `%s`, please wait!", "unexpected status, expected status: `%s`, but the real status: `%s`, please wait!"},
	DockerExecNotFound:      {"`%s`: the docker exec not found", "`%s`: the docker exec not found"},
	DockerImagePullFailed:   {"pull image failed, err: %v", "pull image failed, err: %v"},
	HandlerExecNotFound:     {"`%s`: the handler exec not found", "`%s`: the handler exec not found"},
	CplusActionNotSupport:   {"`%s`: cplus action not support", "`%s`: cplus action not support"},

	ResultUnmarshalFailed:    {"exec result unmarshal failed", "`%s`: exec result unmarshal failed, err: %v"},
	ResultMarshalFailed:      {"exec result marshal failed", "`%s`: exec result marshal failed, err: %v"},
	ChaosbladeServiceStoped:  {"chaosblade service has been stoped", "chaosblade service has been stoped"},
	ProcessIdByNameFailed:    {"system error, uid: `%s`", "`%s`: get process id by name failed, err: %v"},
	ProcessJudgeExistFailed:  {"system error, uid: `%s`", "`%s`: judge the process exist or not, failed, err: %v"},
	ProcessNotExist:          {"system error, uid: `%s`", "`%s`: the process not exist"},
	ProcessGetUsernameFailed: {"system error, uid: `%s`", "`%s`: get username failed by the process id, err: %v"},
	ChannelNil:               {"system error, uid: `%s`", "chanel is nil"},
	SandboxGetPortFailed:     {"system error, uid: `%s`", "get sandbox port failed, err: %v"},
	SandboxCreateTokenFailed: {"system error, uid: `%s`", "create sandbox token failed, err: %v"},
	FileCantGetLogFile:       {"system error, uid: `%s`", "can not get log file"},
	FileNotExist:             {"system error, uid: `%s`", "`%s`: not exist"},
	FileCantReadOrOpen:       {"system error, uid: `%s`", "`%s`: can not read or open"},
	BackfileExists:           {"system error, uid: `%s`", "`%s`: backup file exists, may be annother experiment is running"},
	DbQueryFailed:            {"system error, uid: `%s`", "`%s`: db query failed, err: %v"},
	K8sExecFailed:            {"system error, uid: `%s`", "`%s`: k8s exec failed, err: %v"},
	DockerExecFailed:         {"system error, uid: `%s`", "`%s`: docker exec failed, err: %v"},
	OsCmdExecFailed:          {"system error, uid: `%s`", "`%s`: cmd exec failed, err: %v"},
	HttpExecFailed:           {"system error, uid: `%s`", "`%s`: http cmd failed, err: %v"},
	OsExecutorNotFound:       {"system error, uid: `%s`", "`%s`: os executor not found"},
	ChaosfsClientFailed:      {"init chaosfs client failed in pod %v, err: %v", "init chaosfs client failed in pod %v, err: %v"},
	ChaosfsInjectFailed:      {"inject io exception in pod %s failed, request %v, err: %v", "inject io exception in pod %s failed, request %v, err: %v"},
	ChaosfsRecoverFailed:     {"recover io exception failed in pod  %v, err: %v", "recover io exception failed in pod  %v, err: %v"},
	SshExecFailed:            {"ssh exec failed, result: %v, err %v", "ssh exec failed, result: %v, err %v"},
	SshExecNothing:           {"cannt get reuslt from remote host, please execute recovery and try again", "cannt get reuslt from remote host, please execute recovery and try again"},
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

func ReturnFailWitResult(codeType CodeType, err string, result interface{}) *Response {
	return &Response{Code: codeType.Code, Success: false, Err: err, Result: result}
}

func ReturnSuccess(result interface{}) *Response {
	return &Response{Code: Code[OK].Code, Success: true, Result: result}
}

func ReturnResultIgnoreCode(result interface{}) *Response {
	return &Response{Code: Code[IgnoreCode].Code, Result: result}
}

// new return func for unify errno
func ResponseFailWaitResult(status int32, err string, result interface{}) *Response {
	return &Response{Code: status, Success: false, Err: err, Result: result}
}

func ResponseFail(status int32, err string) *Response {
	return &Response{Code: status, Success: false, Err: err}
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
	err := json.Unmarshal([]byte(content), &resp)
	if err != nil {
		if defaultValue == nil {
			defaultValue = ResponseFailWaitResult(ResultUnmarshalFailed, ResponseErr[ResultUnmarshalFailed].Err,
				fmt.Sprintf(ResponseErr[ResultUnmarshalFailed].ErrInfo, content, err.Error()))
		}
		logrus.Warningf("decode %s err, return default value, %s", content, defaultValue.Print())
		return defaultValue
	}
	return &resp
}
