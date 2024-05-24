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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"log/slog"

	"github.com/PaesslerAG/jsonpath"
)

// used for debug output in this file
const d = true

type Link struct {
	Rel          string `json:"rel,omitempty"`
	Href         string `json:"href,omitempty"`
	Type         string `json:"type,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
}

type Extension struct {
	extentLevel int
	data        any
}

func (p *Processor) loadFile(filename string) (data any, err error) {

	if len(p.inputPath) == 0 {
		p.inputPath = append(p.inputPath, ".")
	}
	for _, path := range p.inputPath {
		testPath := filepath.Join(path, filename)
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			slog.Error(fmt.Sprintf("File %s not found at path %s\n", filename, path))
			continue
		}
		content, err := os.ReadFile(testPath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(content, &data); err != nil {
			slog.Error(fmt.Sprintf("unable to read valid json from %s", testPath), "error", err)
		}
		slog.Info("load file", "path", filename)
		return data, nil
	}
	return nil, fmt.Errorf("file %s not found", filename)
}

var doubleCurlyPattern = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// Process is the main entry point to build a thing description
// out of a thing model, based on the parameters in Processor struct
// but also to process submodel in a top level TM.
func (p *Processor) Process(filename string) error {
	if d {
		slog.Debug("Start Process", "filename", filename, "instance", p.instance.String())
	}
	p.filename = filename
	var err error
	p.data, err = p.loadFile(filename)
	if err != nil {
		return err
	}
	p.iterate(p.data, &PathObject{})
	p.extendAll()
	//copy things to parent
	if p.parent != nil {
		p.copy(p.parent)
	} else {
		p.insertTypeLink()
		thingMap := p.data.(map[string]any)
		thingMap["@type"] = "Thing"
	}
	p.checkVersionInstance()
	if d {
		slog.Debug("End  Process", "filename", filename, "instance", p.instance.String())
	}
	return nil
}

func (p *Processor) copy(to *Processor) {
	p.copyMapSection("properties", to, true)
	p.copyMapSection("actions", to, true)
	p.copyMapSection("events", to, true)
	p.copyPlainSlices("security", to)
	p.copyMapSection("securityDefinitions", to, false)
	p.copyArraySection("links", to)
}

// copyMapSection
func (p *Processor) copyMapSection(section string, to *Processor, usePrefix bool) {
	srcMap := p.data.(map[string]any)
	srcSect, okSrcSect := srcMap[section]
	if okSrcSect {
		destMap := to.data.(map[string]any)
		destSect, okDestSect := destMap[section]
		if !okDestSect {
			destMap[section] = make(map[string]any)
			destSect = destMap[section]
		}
		destSectMap := destSect.(map[string]any)
		srcSectMap := srcSect.(map[string]any)
		prefix := ""
		if usePrefix {
			prefix = p.instance.String()
			if len(prefix) > 0 {
				prefix = prefix + "."
			}
		}
		for k, v := range srcSectMap {
			destSectMap[prefix+k] = v
		}
	}
}

func (p *Processor) copyPlainSlices(section string, to *Processor) {
	srcMap := p.data.(map[string]any)
	srcSect, okSrcSect := srcMap[section]
	if okSrcSect {
		destMap := to.data.(map[string]any)
		destSect, okDestSect := destMap[section]
		if !okDestSect {
			destSect = []any{}
		}
		destSectArray := destSect.([]any)
		srcSectArray := plainOrSliceAsSlice(srcSect)
		destSectArray = append(destSectArray, srcSectArray...)
		destMap[section] = destSectArray
	}
}

func (p *Processor) copyArraySection(section string, to *Processor) {
	srcMap := p.data.(map[string]any)
	srcSect, okSrcSect := srcMap[section]
	if okSrcSect {
		destMap := to.data.(map[string]any)
		destSect, okDestSect := destMap[section]
		if !okDestSect {
			destSect = make(map[string]any)
		}
		destSectArray := destSect.([]any)
		srcSectArray := srcSect.([]any)
		destSectArray = append(destSectArray, srcSectArray...)
		destMap[section] = destSectArray
	}
}

func plainOrSliceAsSlice(in any) []any {
	if plain, ok := in.(string); ok {
		return []any{plain}
	}
	// First, check if it is already string slice
	if out, ok := in.([]any); ok {
		return out
	}
	out := make([]any, len(in.([]any)))
	for i, v := range in.([]any) {
		out[i] = v.(string)
	}
	return out
}

// Save the serialized TD of the already processed TM to
// the defined output
func (p *Processor) Save() {
	// print the result to Outputfile
	prt := NewPrinter()
	printAll(p.data, &PathObject{}, prt, p.VarMap)
	switch p.outputDir {
	case "-":
		fmt.Println(prt.String())
	case "":
	default:
		err := os.MkdirAll(p.outputDir, 0777)
		check(err)
		filename := strings.Replace(p.filename, ".tm.", ".td.", 1)
		err = os.WriteFile(filepath.Join(p.outputDir, filename), prt.ByteArr(), 0644)
		check(err)
	}
}

func (p *Processor) extendAll() {
	for _, e := range p.extensions {
		destMap := p.data.(map[string]any)
		propDest := destMap["properties"].(map[string]any)
		srcMap := e.data.(map[string]any)
		propSrc := srcMap["properties"].(map[string]any)
		merge(propDest, propSrc, 0)
	}
	required, ok := p.data.(map[string]any)["tm:required"]
	delete(p.data.(map[string]any), "tm:required")
	if ok {
		requiredArray := required.([]any)
		for _, requiredElement := range requiredArray {
			requiredString := requiredElement.(string)
			requiredString = strings.Replace(requiredString, "#/", "$.", 1)
			_, notFoundErr := jsonpath.Get(requiredString, p.data)
			if notFoundErr != nil {
				slog.Debug(fmt.Sprintf("Required %s", requiredString), "error", notFoundErr)
			}
		}
	}
}

func (p *Processor) insertTypeLink() {
	destMap := p.data.(map[string]any)
	linksAny, ok := destMap["links"]
	links := linksAny.([]any)
	if !ok {
		links = make([]any, 0, 1)
		destMap["links"] = links
	}
	typelink, _ := structToMap(Link{Rel: "type", Href: p.filename, Type: "application/tm+json"})
	destMap["links"] = append(links, typelink)
}

func merge(dest any, src any, deep int) {
	switch d := dest.(type) {
	case map[string]any:
		srcData, ok := src.(map[string]any)
		if !ok {
			if isInteractiv() {
				fmt.Printf("datatype not a map, skip %T", srcData)
			} else {
				slog.Error(fmt.Sprintf("datatype not a map, skip %T", srcData))
			}
		}
		for key, element := range srcData {
			if dstElement, ok := d[key]; ok {
				merge(dstElement, element, deep+1)
			} else {
				d[key] = element
			}
		}
	case []any:
		// merge of array not yet supported
		//	default:
		//		log.Printf("literal %T %v\n", d, d)
	}

}

func (p *Processor) iterate(data any, po *PathObject) {
	if po.Deep() == 0 {
		slog.Debug(fmt.Sprintf("%siterate %T", indent(po.Deep()), data), "path", po.String(), "deep", po.Deep(), "inst", p.instance.String())
	}
	switch d := data.(type) {
	case map[string]any:
		toDel := make([]string, 0)
		for key, element := range d {
			po.AddMap(key)
			if po.IsPath("links") {
				d[key] = p.processLinks(po, key, element)
			} else if po.IsPath("properties/.*/tm.ref") {
				p.processReference(po, key, element, d)
				toDel = append(toDel, key)
			} else {
				p.iterate(element, po)
			}
			po.Up()
			for _, k := range toDel {
				delete(d, k)
			}
		}
	case []any:
		for i, ele := range d {
			po.AddArray(i)
			p.iterate(ele, po)
			po.Up()
		}
	}
	if po.Deep() == 0 {
		slog.Debug(fmt.Sprintf("%sprocess end  %T", indent(po.Deep()), data), "path", po.String(), "deep", po.Deep())
	}
}

func (p *Processor) processLinks(po *PathObject, key string, element any) []any {
	links := element.([]any)
	returnLinks := make([]any, 0, len(links))
	for _, ele := range links {
		li := ele.(map[string]any)
		if val, ok := li["rel"]; ok && val == "tm:extends" {
			p.foundTMStaff = true
			fileName := li["href"].(string)

			extend, loadError := p.loadFile(fileName)
			p.extensions = append(p.extensions, Extension{extentLevel: po.Deep(), data: extend})
			if loadError != nil {
				slog.Error("unable to read extension", "filname", fileName, "error", loadError)
			}
			p.extensions = append(p.extensions, Extension{extentLevel: po.Deep(), data: extend})
		} else if val, ok := li["rel"]; ok && val == "tm:submodel" {
			p.foundTMStaff = true
			fileName := li["href"].(string)
			pSub := p.NewProcessor()
			if val, ok := li["instanceName"]; ok {
				pSub.instance.AddMap(val.(string))
			} else {
				pSub.instance.AddMap("")
			}
			err := pSub.Process(fileName)
			if err != nil {
				slog.Error("error while processing submodel", "filname", fileName, "error", err)
			}
		} else {
			returnLinks = append(returnLinks, ele)
			p.iterate(ele, po)
		}
	}
	return returnLinks
}

func (p *Processor) processReference(po *PathObject, key string, element any, d map[string]any) {
	ref := strings.Split(element.(string), "#")
	refData, rerr := p.loadFile(ref[0])
	if rerr != nil {
		log.Printf("unable to read reference file: %v\n", rerr)
	} else {
		jpExpr := fmt.Sprintf("$.%s", strings.ReplaceAll(ref[1], "/", "."))
		refDataPart, errp := jsonpath.Get(jpExpr, refData)
		if errp != nil {
			log.Printf("path %s not found in file %s : %v\n", jpExpr, ref[0], errp)
		} else {
			merge(d, refDataPart, po.Deep())
		}
	}
}

func (p *Processor) checkVersionInstance() {
	rootMap := p.data.(map[string]any)
	version, ok := rootMap["version"]
	if ok {
		instVersion := "0.0.0"
		if varsVersion, inVars := p.VarMap["versionInstance"]; inVars {
			instVersion = varsVersion.(string)
		}
		versionMap, okVM := version.(map[string]any)
		if okVM {
			instance, okI := versionMap["instance"]
			if okI {
				instVersion, _ = instance.(string)
			}
		}
		versionMap["instance"] = instVersion
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
