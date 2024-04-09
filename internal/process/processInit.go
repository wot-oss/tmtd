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
	"log"
	"log/slog"
	"strings"
)

type Processor struct {
	extensions   []Extension
	foundTMStaff bool
	VarMap       map[string]any
	//	doubleCurliePattern *regexp.Regexp
	outputDir string
	inputPath []string
	parent    *Processor
	items     []*Processor
	data      any
	filename  string
	instance  PathObject
}

func NewProcessor(out string, in string, vars string) *Processor {
	np := Processor{
		outputDir: out,
		items:     make([]*Processor, 0, 20)}
	np.SetInputPath(in)
	np.SetPlaceholderMap(vars)

	return &np
}

func (p *Processor) NewProcessor() *Processor {
	np := &Processor{outputDir: p.outputDir,
		inputPath: p.inputPath,
		VarMap:    p.VarMap}
	p.items = append(p.items, np)
	np.parent = p
	return np
}

func (p *Processor) SetOutputDir(outputDir string) {
	p.outputDir = outputDir
}

func (p *Processor) SetInputPath(searchPath string) {
	p.inputPath = strings.Split(searchPath, ",")
}

func (p *Processor) SetPlaceholderMap(filename string) {
	if filename != "" {
		varMapAny, err := p.loadFile(filename)
		if err != nil {
			slog.Error("load varMapFile", "filename", filename, "error", err)
			return
		}
		varMap, ok := varMapAny.(map[string]any)
		if !ok {
			log.Printf("varMapFile don't contain a map of values\n")
			return
		}
		p.VarMap = varMap
	}
}

func (p *Processor) String() string {
	return p.instance.String()
}
