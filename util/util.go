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

package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var proPath string
var binPath string
var libPath string
var yamlPath string

// GetProgramPath
func GetProgramPath() string {
	if proPath != "" {
		return proPath
	}
	dir, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal("cannot get the process path")
	}
	if p, err := os.Readlink(dir); err == nil {
		dir = p
	}
	proPath, err = filepath.Abs(filepath.Dir(dir))
	if err != nil {
		log.Fatal("cannot get the full process path")
	}
	return proPath
}

// GetBinPath
func GetBinPath() string {
	if binPath != "" {
		return binPath
	}
	binPath = path.Join(GetProgramPath(), "bin")
	return binPath
}

// GetLibHome
func GetLibHome() string {
	if libPath != "" {
		return libPath
	}
	libPath = path.Join(GetProgramPath(), "lib")
	return libPath
}

func GetYamlHome() string {
	if yamlPath != "" {
		return yamlPath
	}
	yamlPath = path.Join(GetProgramPath(), "yaml")
	return yamlPath
}

// GenerateUid for exp
func GenerateUid() (string, error) {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func IsNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return false
}

//IsExist returns true if file exists
func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// IsDir returns true if the path is directory
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo == nil {
		return false
	}
	return fileInfo.IsDir()
}

//GetUserHome return user home.
func GetUserHome() string {
	user, err := user.Current()
	if err == nil {
		return user.HomeDir
	}
	return "/root"
}

// GetSpecifyingUserHome
func GetSpecifyingUserHome(username string) string {
	usr, err := user.Lookup(username)
	if err == nil {
		return usr.HomeDir
	}
	return fmt.Sprintf("/home/%s", username)
}

// Curl url
func Curl(url string) (string, error, int) {
	logrus.Infoln(url)
	trans := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 10*time.Second)
		},
	}
	client := http.Client{
		Transport: &trans,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err, 0
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err, resp.StatusCode
	}
	return string(bytes), nil, resp.StatusCode
}

// PostCurl
func PostCurl(url string, body []byte, contentType string) (string, error, int) {
	trans := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 10*time.Second)
		},
	}
	client := http.Client{
		Transport: &trans,
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err, 0
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err, 0
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err, response.StatusCode
	}
	return string(bytes), nil, response.StatusCode
}

// CheckPortInUse returns true if the port is in use, otherwise returns false.
func CheckPortInUse(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	if conn != nil {
		return true
	}
	return false
}

func GetUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// GetProgramParentPath returns the parent directory end with /
func GetProgramParentPath() string {
	dir, _ := path.Split(GetProgramPath())
	return dir
}

func GetRunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
