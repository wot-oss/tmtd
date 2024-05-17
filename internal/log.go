/*
Copyright © 2024 Harald Müller <harald.mueller@evosoft.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/innomotics/tmtd/internal/config"
	"github.com/spf13/viper"
)

type DefaultLogHandler struct {
	*slog.TextHandler
}

type DiscardLogHandler struct {
	*slog.TextHandler
}

func newDefaultLogHandler(opts *slog.HandlerOptions) slog.Handler {
	return &DefaultLogHandler{
		TextHandler: slog.NewTextHandler(os.Stderr, opts),
	}
}

func newDiscardLogHandler(opts *slog.HandlerOptions) slog.Handler {
	return &DiscardLogHandler{
		TextHandler: slog.NewTextHandler(io.Discard, opts),
	}
}

func InitLogging() {
	logLevel := viper.GetString(config.KeyLogLevel)

	var logEnabled bool
	level := slog.LevelError

	switch logLevel {
	case "":
		logEnabled = false
	case strings.ToLower(config.LogLevelOff):
		logEnabled = false
	default:
		logEnabled = true
		err := level.UnmarshalText([]byte(logLevel))
		if err != nil {
			level = slog.LevelInfo
		}
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if logEnabled {
		handler = newDefaultLogHandler(opts)
	} else {
		handler = newDiscardLogHandler(opts)
	}

	log := slog.New(handler)
	slog.SetDefault(log)
}
