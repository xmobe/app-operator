// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/xoe-labs/odoo-operator/pkg/controller/odooversion"
)

func main() {
	err := vfsgen.Generate(odooversion.Templates, vfsgen.Options{
		PackageName:  "odooversion",
		BuildTags:    "release",
		VariableName: "Templates",
		Filename:     "zz_generated.templates.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
