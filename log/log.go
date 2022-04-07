package log

import (
	"context"
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/sirupsen/logrus"
	"runtime"
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