/*
 * Copyright 2025 The ChaosBlade Authors
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

package log

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

func Panicf(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Panicf("%s", fmt.Sprintf(format, a...))
}

func Fatalf(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Fatalf("%s", fmt.Sprintf(format, a...))
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Errorf("%s", fmt.Sprintf(format, a...))
}

func Warnf(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Warnf("%s", fmt.Sprintf(format, a...))
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Infof("%s", fmt.Sprintf(format, a...))
}

func Debugf(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Debugf("%s", fmt.Sprintf(format, a...))
}

func Tracef(ctx context.Context, format string, a ...interface{}) {
	uid := ctx.Value(spec.Uid)
	logrus.WithFields(logrus.Fields{
		"uid":      uid,
		"location": GetRunFuncLocation(),
	}).Tracef("%s", fmt.Sprintf(format, a...))
}

func GetRunFuncLocation() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
