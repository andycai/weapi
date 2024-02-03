package activity

import (
	"errors"

	"github.com/andycai/weapi/model"
)

func checkRequest(activityVo *model.Activity) error {
	if activityVo.Name == "" || activityVo.Description == "" {
		return errors.New("name or description required")
	}
	if activityVo.Kind <= 0 || activityVo.Type <= 0 || activityVo.FeeType <= 0 {
		return errors.New("kind, type or fee_type required")
	}
	return nil
}
