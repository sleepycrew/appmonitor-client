package main

import (
	"github.com/a-h/generate"
	"net/url"
	"os"
)

func main() {
	var schemaPath = "./third_party/appmonitor-schema/schema/client.json"
	var schemaUrl, _ = url.Parse("http://appmoniotr.weirdtld/client.json")

	b, err := os.ReadFile(schemaPath) // just pass the file name
	if err != nil {
		panic("couldnt find json schema file.")
	}

	schema := string(b)

	parsedSchema, err := generate.Parse(schema, schemaUrl)
	if err != nil {
		panic("couldnt parse schema!")
	}

	generator := generate.New(parsedSchema)
	generator.CreateTypes()
	if err != nil {
		panic("couldnt create types")
	}

	for i := range generator.Structs {
		gen := generator.Structs[i]
		gen.GenerateCode = false
	}

	w, err := os.Create("./pkg/data/spec.go")
	if err != nil {
		panic("couldnt write spec file")
	}

	generate.Output(w, generator, "data")
}
