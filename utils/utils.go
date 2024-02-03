package utils

import "reflect"

func getType(value interface{}) string {
	t := reflect.TypeOf(value)

	return t.Name()
}

func IsString(value interface{}) bool {
	t := getType(value)
	return t == "string"
}

func IsInt(value interface{}) bool {
	t := getType(value)
	return t == "int"
}

func IsUint(value interface{}) bool {
	t := getType(value)
	return t == "uint"
}

func IsBool(value interface{}) bool {
	t := getType(value)
	return t == "bool"
}
