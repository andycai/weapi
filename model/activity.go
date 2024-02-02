package model

import "time"

type Activity struct {
	ID          uint      `json:"id" gorm:"primarykey;uniqueIndex:,composite:_site_id"`
	SiteID      string    `json:"site_id" gorm:"primaryKey;uniqueIndex:,composite:_site_id"`
	Site        Site      `json:"-"`
	CreatedAt   time.Time `json:"create_at" `
	UpdatedAt   time.Time `json:"update_at" `
	CreatorID   uint      `json:"-"` // 组织者ID
	Creator     User      `json:"-"`
	ClubID      uint      `json:"club_id" gorm:"column:club_id;not null;comment:俱乐部ID"` // 俱乐部ID
	Club        Club      `json:"-"`
	Kind        uint      `json:"kind" gorm:"column:kind;not null;default:1;comment:活动分类:1羽毛球,2篮球,3足球,4聚餐..."`                               // 活动分类:1羽毛球,2篮球,3足球,4聚餐...
	Type        uint      `json:"type" gorm:"column:type;not null;default:1;comment:活动类型:1全局保护,2全局公开,3群组"`                                   // 活动类型:1全局保护,2全局公开,3群组
	Title       string    `json:"title" gorm:"column:title;not null;comment:活动标题"`                                                           // 活动标题
	Description string    `json:"description" gorm:"column:description;not null;comment:活动描述"`                                               // 活动描述
	Quota       uint      `json:"quota" gorm:"column:quota;not null;default:1;comment:报名名额"`                                                 // 报名名额
	Waiting     uint      `json:"waiting" gorm:"column:waiting;not null;default:1;comment:候补数量限制"`                                           // 候补数量限制
	Stage       uint      `json:"stage" gorm:"column:stage;not null;default:1;comment:活动阶段:1报名阶段,2活动阶段,3正常完成和结算,4手动终止活动"`                    // 活动阶段:1报名阶段,2活动阶段,3正常完成和结算,4手动终止活动
	FeeType     uint      `json:"fee_type" gorm:"column:fee_type;not null;default:1;comment:结算方式:1免费,2活动前,3活动后男女平均,4活动后男固定|女平摊,5活动后男平摊|女固定"` // 结算方式:1免费,2活动前,3活动后男女平均,4活动后男固定|女平摊,5活动后男平摊|女固定
	FeeMale     uint      `json:"fee_male" gorm:"column:fee_male;not null;comment:男费用,单位:分"`                                                 // 男费用,单位:分
	FeeFemale   uint      `json:"fee_female" gorm:"column:fee_female;not null;comment:女费用,单位:分"`                                             // 女费用,单位:分
	Addr        string    `json:"addr" gorm:"column:addr;comment:活动地址" json:"addr"`                                                          // 活动地址
	Ahead       uint      `json:"ahead" gorm:"column:ahead;not null;comment:可提前取消时间(小时)"`                                                    // 可提前取消时间(小时)
	BeginAt     time.Time `json:"begin_at" gorm:"column:begin_at;not null;comment:开始时间"`                                                     // 开始时间
	EndAt       time.Time `json:"end_at" gorm:"column:end_at;not null;comment:结束时间"`                                                         // 结束时间
}
