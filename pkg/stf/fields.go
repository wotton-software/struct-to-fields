package stf

import (
	"fmt"
	"reflect"
	"strings"
)

/*
	ExtractFields expects a struct or a pointer to a struct, anything else will result in an error.
	The provided struct instance and its fields->values will be extracted and returned.
*/
func (e *Extractor) ExtractFields(d interface{}) (map[string]interface{}, error) {
	seen := make(map[reflect.Value]bool) //stops code from getting caught in a cycle
	return e.extractFields(d, seen)
}

func (e *Extractor) extractFields(d interface{}, seen map[reflect.Value]bool) (map[string]interface{}, error) {
	elem, err := mustGetStructElem(d)
	if err != nil {
		return nil, err
	}

	n := elem.NumField()

	fields := make(map[string]interface{}, n)
	for i := 0; i < n; i++ {
		newFields, err := e.getFieldsToDataMap(elem, i, seen)
		if err != nil {
			return nil, err
		}

		mergeMaps(fields, newFields)
	}

	return fields, nil
}

func (e *Extractor) getFieldsToDataMap(elem reflect.Value, i int, seen map[reflect.Value]bool) (map[string]interface{}, error) {
	field, tag, fieldName := getFieldDetails(elem, i)
	if e.canIgnore(tag, field) {
		return nil, nil
	}

	if field.CanAddr() {
		if seen[field.Addr()] {
			return nil, nil
		}
		seen[field.Addr()] = true
	}

	fieldName, err := getTagFieldName(fieldName, tag)
	if err != nil {
		return nil, err
	}

	return e.getTags(field, fieldName, seen)
}

func (e *Extractor) canIgnore(tag string, field reflect.Value) bool {
	if hasStfIgnoreTag(tag) {
		return true
	}

	if e.TagRequired && !hasStfTag(tag) {
		return true
	}

	if e.ExcludeNils && isNilable(field.Kind()) && field.IsNil() {
		return true
	}

	return false
}

func (e *Extractor) getTags(field reflect.Value, fieldName string, seen map[reflect.Value]bool) (map[string]interface{}, error) {
	tags := make(map[string]interface{})
	if field.Kind() == reflect.Struct {
		subFields, err := e.extractFields(field.Interface(), seen)
		if err != nil {
			return nil, err
		}

		if len(subFields) == 0 {
			tags[fieldName] = field.Interface()
			return tags, nil
		}

		for k, v := range subFields {
			subName := joinTagNames(fieldName, k)
			tags[subName] = v
		}
	} else {
		tags[fieldName] = field.Interface()
	}

	return tags, nil
}

func isNilable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Ptr:
		return true
	case reflect.Map:
		return true
	case reflect.Array:
		return true
	case reflect.Chan:
		return true
	case reflect.Slice:
		return true
	default:
		return false
	}
}

func mustGetStructElem(d interface{}) (reflect.Value, error) {
	elem := reflect.ValueOf(d)

	originalKind := elem.Kind()

	if originalKind == reflect.Ptr {
		elem = elem.Elem()
	}

	if elem.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("needs struct or pointer to struct type, got: %v", originalKind)
	}

	return elem, nil
}

func mergeMaps(main, sub map[string]interface{}) {
	for k, v := range sub {
		main[k] = v
	}
}

func getField(elem reflect.Value, index int) reflect.Value {
	field := elem.Field(index)
	if field.Kind() == reflect.Ptr {
		f := field.Elem()

		if f != (reflect.Value{}) {
			field = f
		}
	}

	return field
}

func getTagFieldName(fieldName, tagString string) (string, error) {
	switch {
	case hasStfTag(tagString):
		stf, err := getStfTag(tagString)
		if err != nil {
			return "", err
		}

		if !stfTagIsJsonKeyword(stf) {
			fieldName = stf
			break
		}
		fallthrough //json tag means we default to the json:"{value}"
	case hasJsonTag(tagString):
		js, err := getJsonTag(tagString)
		if err != nil {
			return "", err
		}

		fieldName = js
	}

	return fieldName, nil
}

func stfTagIsJsonKeyword(stf string) bool {
	return stf == json
}

func joinTagNames(prefix, suffix string) string {
	return prefix + "." + suffix
}

func getFieldDetails(elem reflect.Value, i int) (reflect.Value, string, string) {
	f := elem.Type().Field(i)
	return getField(elem, i), string(f.Tag), f.Name
}

func trimTagToValue(tag string, prefix string) string {
	tag = strings.ReplaceAll(tag, `"`, "")

	tag = strings.ReplaceAll(tag, prefix, "")

	return tag
}
