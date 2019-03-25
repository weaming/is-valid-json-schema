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
  * `foo/bar?path=a/b/c`: extract data from the POST body by `path` `a/b/c` as the final data to be validate
* `-data`+`-schema`: CLI to verify `data` file against `schema` file
