package converter

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func MapToStruct(mp map[string]interface{}, item interface{}) error {
	for k, v := range mp {
		var jname string
		structValue := reflect.ValueOf(item).Elem()
		fieldByTagName := func(t reflect.StructTag) (string, error) {
			if jt, ok := t.Lookup("keyname"); ok {
				return strings.Split(jt, ",")[0], nil
			}
			return "", fmt.Errorf("tag provided %s does not define a json tag", k)
		}
		fieldNames := map[string]int{}
		for i := 0; i < structValue.NumField(); i++ {
			typeField := structValue.Type().Field(i)
			tag := typeField.Tag
			if string(tag) == "" {
				jname = toMapCase(typeField.Name)
			} else {
				jname, _ = fieldByTagName(tag)
			}
			fieldNames[jname] = i
		}

		fieldNum, ok := fieldNames[k]
		if !ok {
			return fmt.Errorf("field %s does not exist within the provided item", k)
		}
		fieldVal := structValue.Field(fieldNum)
		fieldVal.Set(reflect.ValueOf(v))

	}
	return nil
}

func toMapCase(s string) (str string) {
	runes := []rune(s)
	for j := 0; j < len(runes); j++ {
		if unicode.IsUpper(runes[j]) == true {
			if j == 0 {
				str += strings.ToLower(string(runes[j]))
			} else {
				str += "_" + strings.ToLower(string(runes[j]))
			}
		} else {
			str += strings.ToLower(string(runes[j]))
		}
	}
	return str
}
