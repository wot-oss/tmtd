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

package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	KeyLogLevel             = "logLevel"
	KeyUrlContextRoot       = "urlContextRoot"
	KeyCorsAllowedOrigins   = "corsAllowedOrigins"
	KeyCorsAllowedHeaders   = "corsAllowedHeaders"
	KeyCorsAllowCredentials = "corsAllowCredentials"
	KeyCorsMaxAge           = "corsMaxAge"
	EnvPrefix               = "tmtd"
	LogLevelOff             = "off"
)

var HomeDir string
var DefaultConfigDir string

func InitConfig() {
	var err error
	HomeDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultConfigDir = filepath.Join(HomeDir, ".tmtd")

}

func InitViper() {
	viper.SetDefault("remotes", map[string]any{})
	viper.SetDefault(KeyLogLevel, LogLevelOff)

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(DefaultConfigDir)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; do nothing and rely on defaults
		} else {
			panic("cannot read config: " + err.Error())
		}
	}
	// set prefix "tmtd" for environment variables
	// the environment variables then have to match pattern "tmc_<viper variable>", lower or uppercase
	viper.SetEnvPrefix(EnvPrefix)

	// bind viper variables to environment variables
	_ = viper.BindEnv(KeyLogLevel)             // env variable name = tmtd_loglevel
	_ = viper.BindEnv(KeyUrlContextRoot)       // env variable name = tmtd_urlcontextroot
	_ = viper.BindEnv(KeyCorsAllowedOrigins)   // env variable name = tmtd_corsallowedorigins
	_ = viper.BindEnv(KeyCorsAllowedHeaders)   // env variable name = tmtd_corsallowedheaders
	_ = viper.BindEnv(KeyCorsAllowCredentials) // env variable name = tmtd_corsallowcredentials
	_ = viper.BindEnv(KeyCorsMaxAge)           // env variable name = tmtd_corsmaxage
}
