package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

func ConvertToAVROType(typeString string) string {
	switch typeString {
	case "bool":
		return "boolean"
	}

	return typeString
}

func GenerateAVROSchema[T any](t T) map[string]interface{} {
	var schema = make(map[string]interface{})

	schema["type"] = "struct"
	reflectedType := reflect.TypeOf(t)

	fields := make([]map[string]interface{}, 0)

	// get key of struct T and add it to schema.fields
	for i := 0; i < reflectedType.NumField(); i++ {

		reflectedField := reflectedType.Field(i)
		var field = make(map[string]interface{})

		// get json name of field and add it to schema.fields
		jsonTag := reflectedField.Tag.Get("json")
		if jsonTag != "" {
			// split json tag by comma and get first value
			jsonTag = strings.Split(jsonTag, ",")[0]
			field["field"] = jsonTag
		} else {
			field["field"] = reflectedField.Name
		}

		// get type of field and add it to schema.fields using reflection
		typeOfField := reflectedField.Type

		if typeOfField.Kind() == reflect.Array || typeOfField.Kind() == reflect.Slice {
			field["type"] = "array"
			itemType := typeOfField.Elem()
			if itemType.Kind() == reflect.Struct {
				field["items"] = GenerateAVROSchema(reflect.New(itemType).Elem().Interface())["field"]
			} else {
				field["items"] = ConvertToAVROType(typeOfField.Elem().String())
			}
		} else {
			field["type"] = ConvertToAVROType(typeOfField.String())
		}

		fields = append(fields, field)
	}
	schema["fields"] = fields

	return schema
}

func GenerateAVRO[T any](t T) ([]byte, error) {
	var target = make(map[string]interface{})

	schema := GenerateAVROSchema(t)
	target["schema"] = schema
	target["payload"] = t

	return json.Marshal(target)
}
