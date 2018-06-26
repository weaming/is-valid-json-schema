package main

import (
	"flag"
	"fmt"
	"os"

	fp "path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

var (
	schema string
	data   string
)

func init() {
	flag.StringVar(&schema, "s", "", "schema json path")
	flag.StringVar(&data, "d", "", "data json path")

	flag.Parse()
	if schema == "" || data == "" {
		fmt.Println("Please give the -s -d parameters")
		os.Exit(1)
	}

	var e error
	schema, e = fp.Abs(schema)
	if e != nil {
		fmt.Println(e)
		os.Exit(2)
	}

	data, e = fp.Abs(data)
	if e != nil {
		fmt.Println(e)
		os.Exit(3)
	}
}

func main() {
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%v", schema))
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%v", data))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors:\n")
		for _, desc := range result.Errors() {
			fmt.Printf("* %s\n", desc)
		}
		os.Exit(4)
	}
}
