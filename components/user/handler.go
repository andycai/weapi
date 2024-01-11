package user

import (
	"errors"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

//#region sign in

func SignInPage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)
	if isAuthenticated {
		return c.Redirect("/")
	}

	return core.Render(c, "sign-in/index", fiber.Map{
		"PageTitle":    "Sign In — Werite",
		"FiberCtx":     c,
		"NavBarActive": "sign-in",
	}, "layouts/app")
}

func HTMXSignInPage(c *fiber.Ctx) error {
	return core.Render(c, "sign-in/htmx-sign-in-page", fiber.Map{
		"PageTitle":    "Sign In",
		"NavBarActive": "sign-in",
		"FiberCtx":     c,
	}, "layouts/app-htmx")
}

func HTMXSignInAction(c *fiber.Ctx) error {
	var userVo *model.User
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return core.Render(c, "sign-in/partials/sign-in-form", fiber.Map{
			"Errors": []string{
				"Email or password cannot be null.",
			},
			"IsOob": true,
		}, "layouts/app-htmx")
	}

	err, userVo := Dao.GetByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Render(c, "sign-in/partials/sign-in-form", fiber.Map{
				"Errors": []string{
					"Email and password did not match.",
				},
			}, "layouts/app-htmx")
		}
	}

	if !core.CheckPassword(userVo.Password, password) {
		return core.Render(c, "sign-in/partials/sign-in-form", fiber.Map{
			"Errors": []string{
				"Email and password did not match.",
			},
		}, "layouts/app-htmx")
	}

	Dao.UpdateLoginTime(uint(userVo.ID))
	authentication.AuthStore(c, uint(userVo.ID))

	return core.HTMXRedirectTo("/", "/htmx/home", c)
}

func HTMXSignOut(c *fiber.Ctx) error {
	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/")
	}

	Dao.UpdateLogoutTime(userID)
	authentication.AuthDestroy(c)

	return core.HTMXRedirectTo("/", "/htmx/home", c)
}

//#endregion

//#region sign up

func SignUpPage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)
	if isAuthenticated {
		return c.Redirect("/")
	}

	return core.Render(c, "sign-up/index", fiber.Map{
		"PageTitle":    "Sign Up — Werite",
		"FiberCtx":     c,
		"NavBarActive": "sign-up",
	}, "layouts/app")
}

func HTMXSignUpPage(c *fiber.Ctx) error {
	return core.Render(c, "sign-up/htmx-sign-up-page", fiber.Map{
		"PageTitle":    "Sign Up",
		"NavBarActive": "sign-up",
		"FiberCtx":     c,
	}, "layouts/app-htmx")
}

func HTMXSignUpAction(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// TODO: validate
	if email == "" || username == "" || password == "" {
		return core.Render(c, "sign-up/partials/sign-up-form", fiber.Map{
			"Errors": []string{
				"Username, email, and password cannot be null.",
			},
			"IsOob": true,
		}, "layouts/app-htmx")
	}

	userVo := model.User{Username: username, Email: email, Password: password, Name: username}
	userVo.Password = core.HashPassword(userVo.Password)

	err := Dao.Create(&userVo)

	if err != nil {
		return err
	}

	authentication.AuthStore(c, uint(userVo.ID))

	return core.HTMXRedirectTo("/", "/htmx/home", c)
}

//#endregion
