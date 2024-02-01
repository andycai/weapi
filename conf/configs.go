package conf

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"gorm.io/gorm"
)

// func ReadConf() {
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath(".")

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("Fatal error config file: %s \n", err))
// 	}

// 	pflag.String("app.cacheDir", "./cache/", "cache directory")
// 	pflag.Parse()
// 	viper.BindPFlags(pflag.CommandLine) // bind cmd
// }

var configValueCache *core.ExpiredLRUCache[string, string]

func init() {
	size := 1024 // fixed size
	v, _ := strconv.ParseInt(GetEnv(enum.ENV_CONFIG_CACHE_SIZE), 10, 32)
	if v > 0 {
		size = int(v)
	}

	var configCacheExpired time.Duration = 10 * time.Second
	exp, err := time.ParseDuration(GetEnv(enum.ENV_CONFIG_CACHE_EXPIRED))
	if err == nil {
		configCacheExpired = exp
	}

	configValueCache = core.NewExpiredLRUCache[string, string](size, configCacheExpired)
}

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
	//
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

// load envs to struct
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

func SetValue(db *gorm.DB, key, value string) {
	key = strings.ToUpper(key)
	configValueCache.Remove(key)

	var v model.Config
	result := db.Where("key", key).Take(&v)
	if result.Error != nil {
		newV := &model.Config{
			Key:   key,
			Value: value,
		}
		db.Create(&newV)
		return
	}
	db.Model(&model.Config{}).Where("key", key).UpdateColumn("value", value)
}

func GetValue(db *gorm.DB, key string) string {
	key = strings.ToUpper(key)
	cobj, ok := configValueCache.Get(key)
	if ok {
		return cobj
	}

	var v model.Config
	result := db.Where("key", key).Take(&v)
	if result.Error != nil {
		return ""
	}

	configValueCache.Add(key, v.Value)
	return v.Value
}

func GetIntValue(db *gorm.DB, key string, defaultVal int) int {
	v := GetValue(db, key)
	if v == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultVal
	}
	return int(val)
}

func GetBoolValue(db *gorm.DB, key string) bool {
	v := GetValue(db, key)
	if v == "" {
		return false
	}
	v = strings.ToLower(v)
	if v == "1" || v == "yes" || v == "true" || v == "on" {
		return true
	}
	return false
}

func CheckValue(db *gorm.DB, key, defaultValue string) {
	if GetValue(db, key) == "" {
		SetValue(db, key, defaultValue)
	}
}
