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
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"regexp"
	"strings"
)

// grep ${key}
const ProcessKey = "process"
const ExcludeProcessKey = "excludeProcess"
const ProcessCommandKey = "processCommand"


func GetPidsByLocalPort(ctx context.Context, channel spec.Channel, localPort string) ([]string, error) {
	available := channel.IsCommandAvailable(ctx, "ss")
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
	response := channel.Run(ctx, "ss", fmt.Sprintf("-pln sport = :%s", localPort))
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
	log.Infof(ctx, "sockets for %s, %v", localPort, sockets)
	for idx, s := range sockets {
		if idx == 0 {
			continue
		}
		fields := strings.Fields(s)
		// centos7: users:(("tengine",pid=237768,fd=6),("tengine",pid=237767,fd=6))
		// centos6: users:(("tengine",237768,fd=6),("tengine",237767,fd=6))
		lastField := fields[len(fields)-1]
		log.Infof(ctx,"GetPidsByLocalPort: lastField: %v", lastField)
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
	log.Infof(ctx, "GetPidsByLocalPort: pids: %v", pids)
	return pids, nil
}

func IsAllCommandsAvailable(ctx context.Context, channel spec.Channel, commandNames []string) (*spec.Response, bool) {
	if len(commandNames) == 0 {
		return nil, true
	}

	for _, commandName := range commandNames {
		if channel.IsCommandAvailable(ctx, commandName) {
			continue
		}
		switch commandName {
		case "rm":
			return spec.ResponseFailWithFlags(spec.CommandRmNotFound), false
		case "dd":
			return spec.ResponseFailWithFlags(spec.CommandDdNotFound), false
		case "touch":
			return spec.ResponseFailWithFlags(spec.CommandTouchNotFound), false
		case "mkdir":
			return spec.ResponseFailWithFlags(spec.CommandMkdirNotFound), false
		case "echo":
			return spec.ResponseFailWithFlags(spec.CommandEchoNotFound), false
		case "kill":
			return spec.ResponseFailWithFlags(spec.CommandKillNotFound), false
		case "mv":
			return spec.ResponseFailWithFlags(spec.CommandMvNotFound), false
		case "mount":
			return spec.ResponseFailWithFlags(spec.CommandMountNotFound), false
		case "umount":
			return spec.ResponseFailWithFlags(spec.CommandUmountNotFound), false
		case "tc":
			return spec.ResponseFailWithFlags(spec.CommandTcNotFound), false
		case "head":
			return spec.ResponseFailWithFlags(spec.CommandHeadNotFound), false
		case "grep":
			return spec.ResponseFailWithFlags(spec.CommandGrepNotFound), false
		case "cat":
			return spec.ResponseFailWithFlags(spec.CommandCatNotFound), false
		case "iptables":
			return spec.ResponseFailWithFlags(spec.CommandIptablesNotFound), false
		case "sed":
			return spec.ResponseFailWithFlags(spec.CommandSedNotFound), false
		case "awk":
			return spec.ResponseFailWithFlags(spec.CommandAwkNotFound), false
		case "tar":
			return spec.ResponseFailWithFlags(spec.CommandTarNotFound), false
		}
	}
	return nil, true
}