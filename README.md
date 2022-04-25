# Struct to Fields
A Golang package which converts a given struct or pointer to a struct into a map of field name to interface. 
A typical use case might be when you're performing structured logging and need the json field names.


## Installation
`go get github.com/mikey-wotton/struct-to-fields`

## Examples

* [Logging Use case](https://github.com/wotton-software/struct-to-fields/blob/604c374fa41ee3f97c77f9669a0d24c9ac998296/cmd/logrus-example/main.go#L62)
* [Example output](https://github.com/wotton-software/struct-to-fields/blob/604c374fa41ee3f97c77f9669a0d24c9ac998296/pkg/stf/integration_test.go#L75)

## Additional
### Naming conventions
The naming convention of the extractor uses a hierarchy.
1. If the stf tag is present, use its value as the field name. (myField in example)
2. If the stf tag is not present, use the json tag name. (my_field in example)
3. If neither are present, use the field name. (MyField in example)

```go
type Data struct {
    MyField string `json:"my_field",stf:"myField"`
}
```

### Tag Keywords
Only two STF keywords exist:
 * `stf:"-"` which will stop the extractor from returning this field, useful for sensitive data.
 * `stf:"json"` which will tell the extractor to use the json tag name as the extracted field name. 
 Mostly used when using the tag required option.
```go
type Data struct {
    SensitvieField string `json:"sensitiveField",stf:"-"`
    MyJSONField string `json:"myJsonField",stf:"json"`
}
```
## Options

### Tag Required
By default, the extractor will return all fields within a struct. In logging use cases this 
may be an issue however if you have sensitive data. You can set `TagRequired` using the 
`TagRequiredOption` function when building your extractor. 

```go
extractor := stf.NewExtractor(stf.TagRequiredOption(true))
```


### Exclude Nils
By default, the extractor will return fields which have a nil value. 
You can set `ExcludeNils` using the `ExcludeNilsOption` function when building your extractor. 

```go
extractor := stf.NewExtractor(stf.ExcludeNilsOption(true))
```

