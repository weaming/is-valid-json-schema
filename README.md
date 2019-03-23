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
  * URI endswith `/`: `<uri>/index.json` as schema
  * `foo/bar`: `foo/bar.json` as schema
  * `./map.json` to define your custom mapping from URI to schema file. See [map.json](./map.json)
* `-data`+`-schema`: CLI to varify `schema` file against `data` file
