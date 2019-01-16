// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/xoe-labs/odoo-operator/pkg/controller/odootrack"
)

func main() {
	err := vfsgen.Generate(odootrack.Templates, vfsgen.Options{
		PackageName:  "odootrack",
		BuildTags:    "release",
		VariableName: "Templates",
		Filename:     "zz_generated.templates.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
