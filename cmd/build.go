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
	"fmt"

	"github.com/innomotics/tmtd/internal/process"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "create Thing Descriptions out of models",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("file argument missing")
		}
		//p := &process.Processor{}
		p := process.NewProcessor(cmd.Flag("outputDir").Value.String(),
			cmd.Flag("searchPath").Value.String(),
			cmd.Flag("varmap").Value.String())
		err := p.Process(args[0])
		if err != nil {
			log.Fatal(err)
		}
		p.Save()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("varmap", "m", "", "filename of a json mapfile for substituations")
	buildCmd.Flags().StringP("outputDir", "o", "", "directory for output of thing descriptions")
	buildCmd.Flags().StringP("searchPath", "s", "", "list of directories for source files")

}
