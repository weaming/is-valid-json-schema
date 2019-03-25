package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	libconf "github.com/weaming/golib/config"
	libdj "github.com/weaming/golib/deepjson"
	libfs "github.com/weaming/golib/fs"
	libser "github.com/weaming/golib/serilize"
	gojs "github.com/xeipuuv/gojsonschema"
)

const confPath = "./map.json"

var URI2Path = map[string]string{}

func loadIntoURI2Path() {
	var conf interface{}
	if libfs.IsFile(confPath) {
		libconf.NewConfig("./map.json", &conf)
		for k, v := range conf.(map[string]interface{}) {
			switch v := v.(type) {
			case string:
				URI2Path[k] = v
			default:
				log.Fatal(fmt.Sprintf("%v has type %T\n", v, v))
			}
		}
	}
}

func OK(err error) bool {
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}
	return true
}

func ServeSchemaAsAPI(listen string) {
	http.HandleFunc("/", handler)
	log.Printf("Listening on %v\n", listen)

	// print map
	loadIntoURI2Path()
	bin, _ := libser.JSON(URI2Path)
	fmt.Println(string(bin))

	// serve http
	log.Fatal(http.ListenAndServe(listen, nil))
}

type Response struct {
	Status  int      `json:"status"`
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func NewResponse(status int, success bool, errors []string) Response {
	return Response{status, success, errors}
}

func (r *Response) IsValid() (rv bool) {
	return r.Success == (len(r.Errors) == 0)
}

func getSchemaPath(r *http.Request) string {
	relPath := r.URL.Path
	if strings.HasSuffix(relPath, "/") {
		relPath += "index"
	}

	path := filepath.Join("./", relPath)
	if v, ok := URI2Path[path]; ok {
		path = v
	}
	if !strings.HasSuffix(path, ".json") {
		path += ".json"
	}
	return path
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	schema := getSchemaPath(r)
	log.Printf("%v %v (%v) %v\n", r.Method, r.URL, schema, r.Referer())

	var res Response
	switch r.Method {
	case "GET":
		if !libfs.IsFile(schema) {
			res = NewResponse(404, false, []string{fmt.Sprintf("file %v does not exist", schema)})
		} else {
			w.Write(libfs.ReadFile(schema))
			return
		}

	case "POST":
		data, e := ioutil.ReadAll(r.Body)
		OK(e)

		dataKeyPath := ""
		if v, ok := r.URL.Query()["path"]; ok {
			dataKeyPath = v[0]
		}
		if dataKeyPath != "" {
			dataAny := map[string]interface{}{}
			if e := json.Unmarshal(data, &dataAny); e != nil {
				res = NewResponse(400, false, []string{e.Error()})
				goto writeback
			}
			node := libdj.NewRootNode(dataAny)
			dataTarget := node.GetByPath(strings.Split(dataKeyPath, "/"))
			// rerwite data
			data, e = json.Marshal(dataTarget.Value)
			if e != nil {
				res = NewResponse(400, false, []string{e.Error()})
				goto writeback
			}
		}

		if !libfs.IsFile(schema) {
			res = NewResponse(404, false, []string{fmt.Sprintf("file %v does not exist", schema)})
		} else {
			absSchema, e := filepath.Abs(schema)
			OK(e)
			schemaLoader := gojs.NewReferenceLoader(fmt.Sprintf("file://%v", absSchema))
			documentLoader := gojs.NewStringLoader(string(data))
			valid, err := ValidateSchema(schemaLoader, documentLoader)
			if valid {
				res = NewResponse(200, valid, []string{})
			} else {
				res = NewResponse(400, valid, err.errors)
			}
		}
	default:
		res = NewResponse(405, false, []string{fmt.Sprintf("method %v not allowed", r.Method)})
	}

writeback:
	// write response
	if !res.IsValid() {
		log.Fatal("bug in ServeSchemaAsAPI()")
	}
	bin, err := libser.JSON(res)
	if OK(err) {
		w.WriteHeader(res.Status)
		w.Write(bin)
	} else {
		w.Write([]byte(err.Error()))
	}
}
