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
	"reflect"
	"strings"

	"github.com/mattn/go-isatty"
)

// convert a Struct to a map using the tags
// the lazy way
func structToMap(obj interface{}) (newMap map[string]any, err error) {
	data, err := json.Marshal(obj) // obj -> []byte

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // []byte > map[string]any
	return
}

func indent(indent int) string {
	var sb strings.Builder
	for l := indent; l > 0; l-- {
		sb.WriteRune('\t')
	}
	return sb.String()
}

func isInteractiv() bool {
	fd := os.Stdout.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func logErr(msg string, args ...any) {
	if isInteractiv() {
		fmt.Fprintln(os.Stderr, msg)
	} else {
		slog.Error(msg, args...)
	}
}

func IsNil(val any) bool {
	if val == nil {
		return true
	}
	if reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil() {
		return true
	}
	return false
}
