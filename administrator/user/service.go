package user

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/lib/authentication"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/utils"
	"github.com/gofiber/fiber/v2"
)

var configValueCache *core.ExpiredLRUCache[string, string]

func init() {
	size := 1024 // fixed size
	v, _ := strconv.ParseInt(utils.GetEnv(constant.ENV_CONFIG_CACHE_SIZE), 10, 32)
	if v > 0 {
		size = int(v)
	}

	var configCacheExpired time.Duration = 10 * time.Second
	exp, err := time.ParseDuration(utils.GetEnv(constant.ENV_CONFIG_CACHE_EXPIRED))
	if err == nil {
		configCacheExpired = exp
	}

	configValueCache = core.NewExpiredLRUCache[string, string](size, configCacheExpired)
}

//#region user

func WithAdminAuth(c *fiber.Ctx) error {
	userVo := Current(c)
	signinURL := GetValue(constant.KEY_SITE_SIGNIN_URL)
	if userVo == nil {
		if signinURL == "" {
			return core.Error(c, http.StatusUnauthorized, errors.New("Unauthorized"))
		} else {
			return c.Redirect(signinURL, http.StatusFound)
		}
	}

	if !userVo.IsStaff && !userVo.IsSuperUser {
		return core.Error(c, http.StatusForbidden, errors.New("Forbidden"))
	}

	return c.Next()
}

func SuperAccessCheck(c *fiber.Ctx, obj *model.AdminObject) error {
	if Current(c).IsSuperUser {
		return nil
	}
	return errors.New("only superuser can access")
}

func GetByID(id uint) *model.User {
	var user model.User
	db.Model(&user).
		Where("id = ?", id).
		First(&user)

	return &user
}

func Current(c *fiber.Ctx) *model.User {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = GetByID(userID)
	}

	return userVo
}

func GetByEmail(email string) (error, *model.User) {
	var user model.User
	result := db.Where("email", strings.ToLower(email)).Take(&user)

	return result.Error, &user
}

func Create(user *model.User) error {
	result := db.Create(user)

	return result.Error
}

func UpdatePassword(user *model.User, password string) error {
	p := core.HashPassword(password)
	err := UpdateFields(user, map[string]any{
		"Password": p,
	})
	if err != nil {
		return err
	}
	user.Password = p

	return err
}

func UpdateLogin(c *fiber.Ctx, userID uint) error {
	db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"LastLogin":   time.Now(),
		"LastLoginIP": core.IP(c),
	})

	return nil
}

func UpdateFields(user *model.User, vals map[string]any) error {
	return db.Model(user).Updates(vals).Error
}

//#endregion

//#region config

func CheckConfig() {
	CheckValue(constant.KEY_SITE_LOGO_URL, "/static/img/logo.svg")
	CheckValue(constant.KEY_CMS_GUEST_ACCESS_API, "true")
	CheckValue(constant.KEY_ADMIN_DASHBOARD, "/html/dashboard.html")
	CheckValue(constant.KEY_CMS_UPLOAD_DIR, "./data/uploads/")
	CheckValue(constant.KEY_CMS_MEDIA_PREFIX, "/media/")
	CheckValue(constant.KEY_CMS_MEDIA_HOST, "")
	CheckValue(constant.KEY_CMS_API_HOST, "")
	CheckValue(constant.KEY_CMS_RELATION_COUNT, "3")
	CheckValue(constant.KEY_CMS_SUGGESTION_COUNT, "3")
	CheckValue(constant.KEY_SITE_FAVICON_URL, "/static/img/favicon.png")
	CheckValue(constant.KEY_SITE_SIGNIN_URL, "/auth/login")
	CheckValue(constant.KEY_SITE_SIGNUP_URL, "/auth/register")
	CheckValue(constant.KEY_SITE_LOGOUT_URL, "/auth/logout")
	CheckValue(constant.KEY_SITE_RESET_PASSWORD_URL, "/auth/reset_password")
	CheckValue(constant.KEY_SITE_LOGIN_NEXT, "/")
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

//#endregion
