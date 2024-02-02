package enum

var codeText = map[int]string{
	Success:                       "Success",
	ErrParam:                      "Parameter error",
	ErrData:                       "Data error",
	ErrOp:                         "Operator error",
	ErrTwoPasswordNotMatch:        "Two password is different",
	ErrUserAuth:                   "Authentication failure",
	ErrUserRegister:               "Register error",
	ErrUserGetData:                "Get user data failure",
	ErrUserUpdateData:             "Update user data failure",
	ErrUserNotFound:               "User not found",
	ErrUserEmailFormat:            "Email format is error",
	ErrUserEmailOrPasswordError:   "Email and password did not match",
	ErrUserEmailOrPasswordIsEmpty: "Email or password can not be empty",
	ErrUserDisabled:               "User not allow login",
	ErrUserNotActivated:           "Waiting for activation",
}

func CodeText(code int) string {
	return codeText[code]
}
