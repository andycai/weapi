package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/utils"
)

var configValueCache *core.ExpiredLRUCache[string, string]

func init() {
	size := 1024 // fixed size
	v, _ := strconv.ParseInt(utils.GetEnv(enum.ENV_CONFIG_CACHE_SIZE), 10, 32)
	if v > 0 {
		size = int(v)
	}

	var configCacheExpired time.Duration = 10 * time.Second
	exp, err := time.ParseDuration(utils.GetEnv(enum.ENV_CONFIG_CACHE_EXPIRED))
	if err == nil {
		configCacheExpired = exp
	}

	configValueCache = core.NewExpiredLRUCache[string, string](size, configCacheExpired)
}

func CheckConfig() {
	CheckValue(enum.KEY_SITE_LOGO_URL, "/static/img/logo.svg")
	CheckValue(enum.KEY_CMS_GUEST_ACCESS_API, "true")
	CheckValue(enum.KEY_ADMIN_DASHBOARD, "/html/dashboard.html")
	CheckValue(enum.KEY_CMS_UPLOAD_DIR, "./data/uploads/")
	CheckValue(enum.KEY_CMS_MEDIA_PREFIX, "/media/")
	CheckValue(enum.KEY_CMS_MEDIA_HOST, "")
	CheckValue(enum.KEY_CMS_API_HOST, "")
	CheckValue(enum.KEY_CMS_RELATION_COUNT, "3")
	CheckValue(enum.KEY_CMS_SUGGESTION_COUNT, "3")
	CheckValue(enum.KEY_SITE_FAVICON_URL, "/static/img/favicon.png")
	CheckValue(enum.KEY_SITE_SIGNIN_URL, "/auth/login")
	CheckValue(enum.KEY_SITE_SIGNUP_URL, "/auth/register")
	CheckValue(enum.KEY_SITE_LOGOUT_URL, "/auth/logout")
	CheckValue(enum.KEY_SITE_RESET_PASSWORD_URL, "/auth/reset_password")
	CheckValue(enum.KEY_SITE_LOGIN_NEXT, "/")
}

func SetValue(key, value string) {
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

func GetValue(key string) string {
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

func GetIntValue(key string, defaultVal int) int {
	v := GetValue(key)
	if v == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultVal
	}
	return int(val)
}

func GetBoolValue(key string) bool {
	v := GetValue(key)
	if v == "" {
		return false
	}
	v = strings.ToLower(v)
	if v == "1" || v == "yes" || v == "true" || v == "on" {
		return true
	}
	return false
}

func CheckValue(key, defaultValue string) {
	if GetValue(key) == "" {
		SetValue(key, defaultValue)
	}
}
