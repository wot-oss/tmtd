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

package process

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

func List(searchPath string) {

	inputPath := strings.Split(searchPath, ",")
	if len(inputPath) == 0 {
		inputPath = append(inputPath, "model")
	}
	for _, rootDir := range inputPath {
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, e1 error) error {
			if IsNil(info) {
				logErr("directory did not exist", "path", path)
				return nil
			}
			if !info.IsDir() && strings.HasSuffix(path, ".jsonld") {
				if e1 != nil {
					return e1
				}
				var data any
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(content, &data); err != nil {
					slog.Error("unable to read valid json from stdin", "error", err)
				}
				//atType, _ := jsonpath.Get("$.type", data)
				title, _ := jsonpath.Get("$.title", data)
				fmt.Printf("%60s Title: '%s'\n", path, title)
			}

			return nil
		})
		if err != nil {
			slog.Error("unable to list models dir", "error", err)
		}
	}
}
