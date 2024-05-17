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
	"fmt"
	"regexp"
	"strings"
	//"code.siemens.com/SORIS/soris_main_gobase/util"
)

type PathObject struct {
	path []string
}

func (p *PathObject) AddMap(name string) *PathObject {
	p.path = append(p.path, name)
	return p
}

func (p *PathObject) AddArray(pos int) *PathObject {
	p.path = append(p.path, fmt.Sprintf("[%d]", pos))
	return p
}

func (p *PathObject) AddLiteral() *PathObject {
	p.path = append(p.path, "'literal'")
	return p
}

func (p *PathObject) Deep() int {
	return len(p.path)
}

func (p *PathObject) Up() {
	p.path = p.path[:len(p.path)-1]
}

func (p *PathObject) IsPath(path string) bool {
	pathS := strings.Split(path, "/")
	if len(pathS) != p.Deep() {
		return false
	}
	for i, pat := range pathS {
		reg, err := regexp.Compile(pat)
		if err != nil {
			fmt.Printf("path pattern: %v", err)
			return false
		}
		if !reg.MatchString(p.path[i]) {
			return false
		}
	}
	return true
}

func (p *PathObject) Prefix() string {
	if len(p.path) > 1 && p.path[0] == "properties" {
		parts := strings.Split(p.path[1], ".")
		return strings.Join(parts[:len(parts)-1], ".")
	}
	return ""
}

func (p *PathObject) String() string {
	var b strings.Builder
	sep := ""
	for _, part := range p.path {
		fmt.Fprintf(&b, "%s%s", sep, part)
		sep = "/"
	}
	return b.String()
}
