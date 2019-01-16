// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/xoe-labs/odoo-operator/pkg/controller/odooinstance"
)

func main() {
	err := vfsgen.Generate(odooinstance.Templates, vfsgen.Options{
		PackageName:  "odooinstance",
		BuildTags:    "release",
		VariableName: "Templates",
		Filename:     "zz_generated.templates.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
