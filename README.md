# tmtd

[![Go Report Card](https://goreportcard.com/badge/github.com/innomotics/tmtd)](https://goreportcard.com/report/github.com/innomotics/tmtd) [![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/wot-oss/tmtd)](https://github.com/innomotics/tmtd/releases) [![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/innomotics/tmtd)

Transpiling thing-model to thing-descriptions

⚠ This software is **experimental** and may not be fit for any purposes. 

# examples

### use tm:extend 
tmtd build -m vars.json -o thing -s model dim.jsonld

### use tm:submodel 
tmtd build -m vars.json -o thing -s model SmartVentilator.tm.jsonld

### use a more complex model from W3C 
tmtd build -m vars.json -o thing -s model/w3cTest floor-lamp-1.0.0.tm.jsonld
