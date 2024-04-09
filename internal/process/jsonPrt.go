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
	"log"
	"log/slog"
	"slices"
	"strings"
)

type Printer struct {
	res []byte
}

func NewPrinter() *Printer {
	p := Printer{res: make([]byte, 0, 5000)}
	return &p
}

func (pr *Printer) Add(c any) {
	b, err := json.Marshal(c)
	if err != nil {
		slog.Error("try to add", "error", err)
	}
	//slog.Debug(fmt.Sprintf("Add to output: '%v'", c))
	pr.res = append(pr.res, b...)
}

func (pr *Printer) AddText(s string) {
	b := []byte(s)
	//log.Printf("Add to output: '%v'\n", s)
	pr.res = append(pr.res, b...)
}

func (pr *Printer) AddInterface(found any, isString bool) {
	pr.Add(found)
}

func (pr *Printer) String() string {
	return string(pr.res)
}

func (pr *Printer) ByteArr() []byte {
	return pr.res
}

func (pr *Printer) addNl(indent int) {
	pr.AddText("\n")
	for l := indent; l > 0; l-- {
		pr.AddText("\t")
	}
}

func (pr *Printer) addSep(needSep bool, indent int) bool {
	if needSep {
		pr.AddText(",")
		pr.addNl(indent)
	}
	return true
}

func printAll(data any, po *PathObject, jsPrt *Printer, vars map[string]any) {
	var needSep bool
	slog.Debug(fmt.Sprintf("print %T %v\n", data, data), "path", po.String(), "deep", po.Deep())

	order := []string{"@context", "title", "@type", "description", "version", "securityDefinitions", "security", "links", "properties", "actions", "events"}

	topMap, ok := data.(map[string]any)
	for k := range topMap {
		if !slices.Contains(order, k) {
			order = append(order, k)
		}
	}
	if ok {
		jsPrt.AddText("{")
		jsPrt.addNl(po.Deep() + 1)
		for _, key := range order {
			var pd any
			pd, ok := topMap[key]
			if ok && pd != nil {
				po.AddMap(key)
				needSep = jsPrt.addSep(needSep, po.Deep())
				jsPrt.AddText(fmt.Sprintf("\"%s\": ", key))
				printRest(pd, po, jsPrt, vars)
				po.Up()
			}
		}
		jsPrt.addNl(po.Deep())
		jsPrt.AddText("}")
	}

}

func printRest(data any, po *PathObject, jsPrt *Printer, vars map[string]any) {
	//	slog.Debug(fmt.Sprintf("print %T %v\n", data, data), "path", po.String(), "deep", po.Deep())
	//indent := true
	var needSep bool
	switch d := data.(type) {
	case map[string]any:
		jsPrt.AddText("{")
		if len(d) > 0 {
			jsPrt.addNl(po.Deep() + 1)
		}
		needSep = false
		// create a sorted list of keys
		keys := make([]string, 0, len(d))
		for k := range d {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, key := range keys {
			element := d[key]
			po.AddMap(key)
			needSep = jsPrt.addSep(needSep, po.Deep())
			jsPrt.AddText(fmt.Sprintf("\"%s\": ", key))
			printRest(element, po, jsPrt, vars)
			po.Up()
		}
		if len(d) > 0 {
			jsPrt.addNl(po.Deep())
		}
		jsPrt.AddText("}")
	case ([]any):
		jsPrt.AddText("[")
		needSep = false
		for i, ele := range d {
			po.AddArray(i)
			needSep = jsPrt.addSep(needSep, po.Deep())
			printRest(ele, po, jsPrt, vars)
			po.Up()
		}
		if len(d) > 0 {
			jsPrt.addNl(po.Deep())
		}
		jsPrt.AddText("]")
	case string:
		varName := strings.Trim(d, " \t")
		if strings.HasPrefix(varName, "{{") && strings.HasSuffix(varName, "}}") {
			res := doubleCurlyPattern.FindStringSubmatch(d)
			if len(res) == 2 {
				found, ok := vars[res[1]]
				_, isString := found.(string)
				if ok {
					jsPrt.AddInterface(found, isString)
					break
				}
			}
		}
		res := doubleCurlyPattern.ReplaceAllStringFunc(d, func(m string) string {
			varName := strings.Replace(m, "{{", "", 1)
			varName = strings.Replace(varName, "}}", "", 1)
			varName = strings.Trim(varName, " \t")
			found, ok := vars[varName]
			if ok {
				return fmt.Sprintf("%v", found)
			}
			return m
		})
		ok := true
		if ok {
			jsPrt.AddInterface(res, ok)
		} else {
			jsPrt.AddText(d)
		}

	case bool:
		jsPrt.Add(d)
	case json.Number, float32, float64, int, int16, int32, int64, int8, uint:

		jsPrt.AddInterface(d, false)
	default:
		log.Printf("literal %T %v\n", d, d)
	}
}
