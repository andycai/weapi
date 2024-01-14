package enum

import "errors"

const (
	ErrorTextGroupNotFound        = "群组不存在"
	ErrorTextGroupFireMember      = "只有群管理才能踢人出群"
	ErrorTextGroupRemove          = "只有群主才能删除群组"
	ErrorTextGroupPromote         = "只有群主才能提升管理员"
	ErrorTextGroupTransfer        = "只有群主才能转让群主职位"
	ErrorTextGroupAlreadyManager  = "已经时管理员，不需要再提升"
	ErrorTextGroupOwnerCannotQuit = "群主不能退群，请先转让群主后再退群"
	ErrorTextGroupManagerFull     = "群管理员名额已满"
	ErrorTextGroupMemberFull      = "群成员名额已满"
	ErrorTextGroupNameExists      = "群的名字已经存在"
	ErrorTextGroupManagerOp       = "只有群管理员才能操作"
	ErrorTextActivityNotFound     = "活动不存在"
)

var errorDict map[string]error

func GetError(text string) error {
	if errorDict == nil {
		errorDict = make(map[string]error)
	}

	if _, ok := errorDict[text]; !ok {
		errorDict[text] = errors.New(text)
	}

	return errorDict[text]
}
