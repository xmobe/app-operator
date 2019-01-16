// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/xoe-labs/odoo-operator/pkg/controller/odoocluster"
)

func main() {
	err := vfsgen.Generate(odoocluster.Templates, vfsgen.Options{
		PackageName:  "odoocluster",
		BuildTags:    "release",
		VariableName: "Templates",
		Filename:     "zz_generated.templates.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
