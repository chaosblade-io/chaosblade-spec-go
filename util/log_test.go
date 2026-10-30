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

package util

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestInitLogConsumesLogLevel verifies InitLog resolves the effective level from
// Debug / LogLevel / the CHAOSBLADE_LOG_LEVEL environment variable and exports it
// for spawned binaries. InitLog mutates process-global state, so it is saved and
// restored to keep the test hermetic.
func TestInitLogConsumesLogLevel(t *testing.T) {
	origDebug, origLogLevel, origLogPath := Debug, LogLevel, LogPath
	origLevel := logrus.GetLevel()
	origEnv, hadEnv := os.LookupEnv(LogLevelEnv)
	defer func() {
		Debug, LogLevel, LogPath = origDebug, origLogLevel, origLogPath
		logrus.SetLevel(origLevel)
		if hadEnv {
			os.Setenv(LogLevelEnv, origEnv)
		} else {
			os.Unsetenv(LogLevelEnv)
		}
	}()

	t.Run("blade honors flag and ignores env", func(t *testing.T) {
		Debug = false
		LogLevel = "error"
		// A stray env var in the user's shell must not override the CLI flag.
		os.Setenv(LogLevelEnv, "debug")
		logrus.SetLevel(logrus.InfoLevel)
		InitLog(Blade)
		if got := logrus.GetLevel(); got != logrus.ErrorLevel {
			t.Errorf("Blade level = %v, want error", got)
		}
		if got := os.Getenv(LogLevelEnv); got != "error" {
			t.Errorf("exported %s = %q, want error", LogLevelEnv, got)
		}
	})

	t.Run("debug wins over log-level", func(t *testing.T) {
		Debug = true
		LogLevel = "error"
		logrus.SetLevel(logrus.InfoLevel)
		InitLog(Blade)
		if got := logrus.GetLevel(); got != logrus.DebugLevel {
			t.Errorf("level = %v, want debug", got)
		}
		Debug = false
	})

	t.Run("spawned binary inherits exported level", func(t *testing.T) {
		Debug = false
		LogLevel = "info" // the default a spawned binary is compiled with
		LogPath = t.TempDir()
		os.Setenv(LogLevelEnv, "warn")
		logrus.SetLevel(logrus.InfoLevel)
		InitLog(Custom) // programType != Blade -> env takes precedence
		if got := logrus.GetLevel(); got != logrus.WarnLevel {
			t.Errorf("inherited level = %v, want warn", got)
		}
	})

	t.Run("invalid level falls back without panic", func(t *testing.T) {
		Debug = false
		LogLevel = "bogus"
		os.Unsetenv(LogLevelEnv)
		logrus.SetLevel(logrus.InfoLevel)
		InitLog(Blade)
		if got := logrus.GetLevel(); got != logrus.InfoLevel {
			t.Errorf("level = %v, want info", got)
		}
	})
}
