package model

import "time"

const (
	ActivityKindBasketball = iota + 1
	ActivityKindFootball
	ActivityKindVolleyball
	ActivityKindBadminton
	ActivityKindDinner
	ActivityKindOther
)

const (
	ActivityTypePrivate = iota + 1
	ActivityTypePublic
	ActivityTypeClub
)

const (
	ActivityStageSignup = iota + 1
	ActivityStageInProcess
	ActivityStageFinished
	ActivityStageAbort
)

const (
	ActivityFeeTypeFree = iota + 1
	ActivityFeeTypeFixed
	ActivityFeeTypeAA
	ActivityFeeTypeMaleFixedFemaleAA
	ActivityFeeTypeMaleAAFemaleFixed
)

type Activity struct {
	BaseField
	ClubID    uint      `json:"club_id" gorm:"index"`
	Club      Club      `json:"-"`
	Kind      uint      `json:"kind" gorm:"index;default:1;comment:活动分类:1羽毛球,2篮球,3足球,4聚餐..."`                     // 活动分类:1羽毛球,2篮球,3足球,4聚餐...
	Type      uint      `json:"type" gorm:"index;default:1;comment:活动类型:1全局保护,2全局公开,3群组"`                         // 活动类型:1全局保护,2全局公开,3群组
	Quota     uint      `json:"quota" gorm:"default:1;comment:报名名额"`                                              // 报名名额
	Waiting   uint      `json:"waiting" gorm:"default:1;comment:候补数量限制"`                                          // 候补数量限制
	Stage     uint      `json:"stage" gorm:"default:1;comment:活动阶段:1报名阶段,2活动阶段,3正常完成和结算,4手动终止活动"`                 // 活动阶段:1报名阶段,2活动阶段,3正常完成和结算,4手动终止活动
	FeeType   uint      `json:"fee_type" gorm:"default:1;comment:结算方式:1免费,2活动前,3活动后男女平均,4活动后男固定|女平摊,5活动后男平摊|女固定"` // 结算方式:1免费,2活动前,3活动后男女平均,4活动后男固定|女平摊,5活动后男平摊|女固定
	FeeMale   uint      `json:"fee_male" gorm:"comment:男费用,单位:分"`                                                 // 男费用,单位:分
	FeeFemale uint      `json:"fee_female" gorm:"comment:女费用,单位:分"`                                               // 女费用,单位:分
	Address   string    `json:"address" gorm:"comment:活动地址"`                                                      // 活动地址
	Ahead     uint      `json:"ahead" gorm:"default:24;comment:可提前取消时间(小时)"`                                      // 可提前取消时间(小时)
	BeginAt   time.Time `json:"begin_at" gorm:"index;comment:开始时间"`                                               // 开始时间
	EndAt     time.Time `json:"end_at" gorm:"index;comment:结束时间"`                                                 // 结束时间
}
