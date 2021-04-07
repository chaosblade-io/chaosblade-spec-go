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
	"flag"
	"io"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	Blade = 1
	Bin   = 2
)

const BladeLog = "chaosblade.log"

var (
	Debug    bool
	LogLevel string
)

func AddDebugFlag() {
	flag.BoolVar(&Debug, "debug", false, "set debug mode")
}

func AddLogLevelFlag() {
	flag.StringVar(&LogLevel, "log-level", "info", "level of logging wanted.")
}

// InitLog invoked after flag parsed
func InitLog(programType int) {
	logFile, err := GetLogFile(programType)
	if err != nil {
		return
	}
	output := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    30, // m
		MaxBackups: 1,
		MaxAge:     2, // days
		Compress:   false,
	}
	logrus.SetOutput(&fileWriterWithoutErr{output})

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	}
	logrus.SetFormatter(formatter)

	if Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func GetLogPath(programType int) (string, error) {
	var binDir string
	switch programType {
	case Blade:
		binDir = GetProgramPath()
	case Bin:
		binDir = GetProgramParentPath()
	default:
		binDir = GetProgramPath()
	}
	logsPath := path.Join(binDir, "logs")
	if IsExist(logsPath) {
		return logsPath, nil
	}
	// mk dir
	err := os.MkdirAll(logsPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return logsPath, nil
}

// GetLogFile
func GetLogFile(programType int) (string, error) {
	logPath, err := GetLogPath(programType)
	if err != nil {
		return "", err
	}
	logFile := path.Join(logPath, BladeLog)
	return logFile, nil
}

// GetNohupOutput
func GetNohupOutput(programType int, logFileName string) string {
	logPath, err := GetLogPath(programType)
	if err != nil {
		return "/dev/null"
	}
	return path.Join(logPath, logFileName)
}

// fileWriterWithoutErr write func does not return err under any conditions
// To solve "Failed to write to log, write logs/chaosblade.log: no space left on device" err
type fileWriterWithoutErr struct {
	io.Writer
}

func (f *fileWriterWithoutErr) Write(b []byte) (n int, err error) {
	i, _ := f.Writer.Write(b)
	return i, nil
}
