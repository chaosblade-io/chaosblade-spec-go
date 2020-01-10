package util

import (
	"fmt"
	"io"
	syslog "log"
	"os"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type encoderConfigFunc func(*zapcore.EncoderConfig)

type config struct {
	encoder zapcore.Encoder
	level   zap.AtomicLevel
	opts    []zap.Option
}

func Logger(output io.Writer) logr.Logger {
	return LoggerTo(output)
}

func LoggerTo(destWriter io.Writer) logr.Logger {
	conf := getConfig()
	return createLogger(conf, destWriter)
}

func createLogger(conf config, destWriter io.Writer) logr.Logger {
	syncer := zapcore.AddSync(destWriter)
	conf.encoder = &logf.KubeAwareEncoder{Encoder: conf.encoder, Verbose: conf.level.Level() < 0}
	conf.opts = append(conf.opts, zap.AddCallerSkip(1), zap.ErrorOutput(syncer))
	multSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), syncer)
	log := zap.New(zapcore.NewCore(conf.encoder, multSyncer, conf.level))
	log = log.WithOptions(conf.opts...)
	return zapr.NewLogger(log)
}

func getConfig() config {
	var c config

	var newEncoder func(...encoderConfigFunc) zapcore.Encoder

	// Set the defaults depending on the log mode (development vs. production)
	if Debug {
		newEncoder = newConsoleEncoder
		c.level = zap.NewAtomicLevelAt(zap.DebugLevel)
		c.opts = append(c.opts, zap.Development(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		newEncoder = newJSONEncoder
		c.level = zap.NewAtomicLevelAt(zap.InfoLevel)
		c.opts = append(c.opts, zap.AddStacktrace(zap.WarnLevel))
	}

	var ecfs []encoderConfigFunc

	ecfs = append(ecfs, withTimeEncoding(zapcore.RFC3339NanoTimeEncoder))

	c.encoder = newEncoder(ecfs...)

	err := c.setLogLevel(LogLevel)
	if err != nil {
		syslog.Fatal(err)
	}
	return c
}

func newJSONEncoder(ecfs ...encoderConfigFunc) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	for _, f := range ecfs {
		f(&encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func newConsoleEncoder(ecfs ...encoderConfigFunc) zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	for _, f := range ecfs {
		f(&encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func withTimeEncoding(te zapcore.TimeEncoder) encoderConfigFunc {
	return func(ec *zapcore.EncoderConfig) {
		ec.EncodeTime = te
	}
}

func (c *config) setLogLevel(l string) error {
	lower := strings.ToLower(l)
	var lvl int
	switch lower {
	case "debug":
		lvl = -1
	case "info":
		lvl = 0
	case "error":
		lvl = 2
	default:
		i, err := strconv.Atoi(lower)
		if err != nil {
			return fmt.Errorf("invalid log level \"%s\"", l)
		}

		if i > 0 {
			lvl = -1 * i
		} else {
			return fmt.Errorf("invalid log level \"%s\"", l)
		}
	}
	level := zapcore.Level(int8(lvl))
	c.level = zap.NewAtomicLevelAt(level)
	return nil
}
