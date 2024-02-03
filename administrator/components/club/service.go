package club

import (
	"errors"

	"github.com/andycai/weapi/model"
)

func checkRequest(clubVo *model.Club) error {
	if clubVo.Name == "" || clubVo.Description == "" {
		return errors.New("name or description required")
	}
	return nil
}

func checkMemberRequest(clubMemberVo *model.ClubMember) error {
	if clubMemberVo.DisplayName == "" {
		return errors.New("display name required")
	}
	if clubMemberVo.Position <= 0 {
		return errors.New("position  required")
	}
	return nil
}
