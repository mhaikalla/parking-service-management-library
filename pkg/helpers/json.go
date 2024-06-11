package helpers

import "reflect"

func GetAllFieldObj(obj interface{}) []string {
	userType := reflect.TypeOf(obj)
	jsonFields := make([]string, 0)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			jsonFields = append(jsonFields, jsonTag)
		}
	}
	return jsonFields
}
