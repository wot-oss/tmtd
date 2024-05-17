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
package cmd

import (
	"log/slog"
	"slices"

	"github.com/innomotics/tmtd/internal"
	"github.com/innomotics/tmtd/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tmtd",
	Short: "Transpiling a thing-model to a concrete thing-description",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var loglevel string
var logEnabledDefaultCmd = []string{"serve"}

func init() {

	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "", "enable logging by setting a log level, one of [error, warn, info, debug, off]")
	rootCmd.PersistentPreRun = preRunAll

	config.InitViper()
	// bind viper variable "loglevel" to CLI flag --loglevel of root command
	_ = viper.BindPFlag(config.KeyLogLevel, rootCmd.PersistentFlags().Lookup("loglevel"))
}

func preRunAll(cmd *cobra.Command, args []string) {
	// set default loglevel depending on subcommand
	logDefault := cmd != nil && slices.Contains(logEnabledDefaultCmd, cmd.CalledAs())
	if logDefault {
		viper.SetDefault(config.KeyLogLevel, slog.LevelInfo.String())
	} else {
		viper.SetDefault(config.KeyLogLevel, config.LogLevelOff)
	}

	internal.InitLogging()
}
