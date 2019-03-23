package main

import (
	"flag"
	"fmt"
	"os"
	fp "path/filepath"

	gojc "github.com/xeipuuv/gojsonschema"
)

var (
	schema string
	data   string
	listen string
)

func parseArgs() bool {
	flag.StringVar(&schema, "schema", "", "schema json path")
	flag.StringVar(&data, "data", "", "data json path")
	flag.StringVar(&listen, "listen", "", "fake api server listen on")

	flag.Parse()
	if listen == "" {
		if schema == "" || data == "" {
			fmt.Println("Please give the -schema -data parameters")
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
		return true
	} else {
		return false
	}
}

func main() {
	if parseArgs() {
		schemaLoader := gojc.NewReferenceLoader(fmt.Sprintf("file://%v", schema))
		documentLoader := gojc.NewReferenceLoader(fmt.Sprintf("file://%v", data))
		valid, err := ValidateSchema(schemaLoader, documentLoader)
		if valid {
			fmt.Println("The document is valid")
		} else {
			fmt.Printf("The document is not valid. see errors:\n")
			fmt.Println(err.Error())
			os.Exit(4)
		}
	} else {
		ServeSchemaAsAPI(listen)
	}
}
