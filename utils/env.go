package utils

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

func GetEnv(key string) string {
	v, _ := LookupEnv(key)
	return v
}

func GetIntEnv(key string) int {
	v := GetEnv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		return i
	}
	return 0
}

func GetBoolEnv(key string) bool {
	v := strings.ToLower(GetEnv(key))
	return v == "1" || v == "yes" || v == "true" || v == "on"
}

func LookupEnv(key string) (string, bool) {
	// Check .env file
	data, err := os.ReadFile(".env")
	if err != nil {
		return os.LookupEnv(key)
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines); i++ {
		v := strings.TrimSpace(lines[i])
		if v == "" {
			continue
		}
		if v[0] == '#' {
			continue
		}
		if !strings.Contains(v, "=") {
			continue
		}
		vs := strings.SplitN(v, "=", 2)
		if strings.EqualFold(strings.TrimSpace(vs[0]), key) {
			return strings.TrimSpace(vs[1]), true
		}
	}
	return "", false
}

// LoadEnvs load envs to struct
func LoadEnvs(objPtr any) {
	if objPtr == nil {
		return
	}
	elm := reflect.ValueOf(objPtr).Elem()
	elmType := elm.Type()

	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		if !f.CanSet() {
			continue
		}
		keyName := elmType.Field(i).Tag.Get("env")
		if keyName == "-" {
			continue
		}
		if keyName == "" {
			keyName = elmType.Field(i).Name
		}
		switch f.Kind() {
		case reflect.String:
			if v, ok := LookupEnv(keyName); ok {
				f.SetString(v)
			}
		case reflect.Int:
			if v, ok := LookupEnv(keyName); ok {
				if iv, err := strconv.ParseInt(v, 10, 32); err == nil {
					f.SetInt(iv)
				}
			}
		case reflect.Bool:
			if v, ok := LookupEnv(keyName); ok {
				v := strings.ToLower(v)
				yes := v == "1" || v == "yes" || v == "true" || v == "on"
				f.SetBool(yes)
			}
		}
	}
}
