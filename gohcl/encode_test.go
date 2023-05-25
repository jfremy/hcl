// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gohcl_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ExampleEncodeIntoBody() {
	type Service struct {
		Name string   `hcl:"name,label"`
		Exe  []string `hcl:"executable"`
	}
	type Constraints struct {
		OS   string `hcl:"os"`
		Arch string `hcl:"arch"`
	}
	type App struct {
		Name        string       `hcl:"name"`
		Desc        string       `hcl:"description"`
		Constraints *Constraints `hcl:"constraints,block"`
		Services    []Service    `hcl:"service,block"`
	}

	app := App{
		Name: "awesome-app",
		Desc: "Such an awesome application",
		Constraints: &Constraints{
			OS:   "linux",
			Arch: "amd64",
		},
		Services: []Service{
			{
				Name: "web",
				Exe:  []string{"./web", "--listen=:8080"},
			},
			{
				Name: "worker",
				Exe:  []string{"./worker"},
			},
		},
	}

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&app, f.Body())
	fmt.Printf("%s", f.Bytes())

	// Output:
	// name        = "awesome-app"
	// description = "Such an awesome application"
	//
	// constraints {
	//   os   = "linux"
	//   arch = "amd64"
	// }
	//
	// service "web" {
	//   executable = ["./web", "--listen=:8080"]
	// }
	// service "worker" {
	//   executable = ["./worker"]
	// }
}

func TestNonRegressionBlockMap(t *testing.T) {
	type Child struct {
		Name     string `hcl:"name,label"`
		Property string `hcl:"prop,optional"`
	}
	type Base struct {
		TestMap map[string]*Child      `hcl:"testMap,block"`
		Env     map[string]string      `hcl:"env,block"`
		Config  map[string]interface{} `hcl:"config,block"`
	}

	base := Base{
		TestMap: map[string]*Child{
			"test": {Name: "test", Property: "property"},
		},
		Env: map[string]string{
			"PATH": "/bin:/sbin",
		},
		Config: map[string]interface{}{
			"cmd":   []string{"/bin/ls"},
			"image": "busybox",
		},
	}

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&base, f.Body())
	fmt.Printf("%s", f.Bytes())
}
