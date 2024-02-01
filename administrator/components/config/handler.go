package config

import (
	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/enum"
)

func CheckConfig() {
	conf.CheckValue(db, enum.KEY_SITE_LOGO_URL, "/static/img/logo.svg")
	conf.CheckValue(db, enum.KEY_CMS_GUEST_ACCESS_API, "true")
	conf.CheckValue(db, enum.KEY_ADMIN_DASHBOARD, "./dashboard.html")
	conf.CheckValue(db, enum.KEY_CMS_UPLOAD_DIR, "./data/uploads/")
	conf.CheckValue(db, enum.KEY_CMS_MEDIA_PREFIX, "/media/")
	conf.CheckValue(db, enum.KEY_CMS_MEDIA_HOST, "")
	conf.CheckValue(db, enum.KEY_CMS_API_HOST, "")
	conf.CheckValue(db, enum.KEY_CMS_RELATION_COUNT, "3")
	conf.CheckValue(db, enum.KEY_CMS_SUGGESTION_COUNT, "3")
	conf.CheckValue(db, enum.KEY_SITE_FAVICON_URL, "/static/img/favicon.png")
	conf.CheckValue(db, enum.KEY_SITE_SIGNIN_URL, "/auth/login")
	conf.CheckValue(db, enum.KEY_SITE_SIGNUP_URL, "/auth/register")
	conf.CheckValue(db, enum.KEY_SITE_LOGOUT_URL, "/auth/logout")
	conf.CheckValue(db, enum.KEY_SITE_RESET_PASSWORD_URL, "/auth/reset_password")
	conf.CheckValue(db, enum.KEY_SITE_LOGIN_NEXT, "/")
}
