package enum

var codeText = map[int]string{
	Success:                       "Success",
	SucGroupApply:                 "申请加入群组成功！",
	ErrParam:                      "Parameter error.",
	ErrData:                       "Data error.",
	ErrOp:                         "Operator error.",
	ErrTwoPasswordNotMatch:        "输入的两次密码不一致！",
	ErrUserAuth:                   "登录验证失败，请重新登录！",
	ErrUserRegister:               "注册失败！",
	ErrUserData:                   "获取用户数据失败！",
	ErrUserUpdateData:             "更新用户数据失败！",
	ErrUserNotFound:               "用户不存在！",
	ErrUserEmailFormat:            "Email format is error.",
	ErrUserEmailOrPasswordError:   "Email and password did not match.",
	ErrUserEmailOrPasswordIsEmpty: "Email or password can not be empty.",
}

func CodeText(code int) string {
	return codeText[code]
}
