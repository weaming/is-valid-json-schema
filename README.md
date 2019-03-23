# isFuckingValidJSONSchema
## Install

```
go get -u github.com/weaming/isFuckingValidJSONSchema
```

## Usage

```
Usage of isFuckingValidJSONSchema:
  -data string
    	data json path
  -listen string
    	fake api server listen on
  -schema string
    	schema json path
```

* `-listen`: serve `$CWD` filesystem as API to provide JSON schema validation service.
  * URI endswith '/': `<uri>/index.jon` as schema
  * `foo/bar`: `foo/bar.jon` as schema
* `-data`+`-schema`: CLI to varify JSON file `schema` against JSON file `data`
