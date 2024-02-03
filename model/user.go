package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/andycai/weapi/utils/random"
	"gorm.io/gorm"
)

const (
	PermissionAll    = "all"
	PermissionCreate = "create"
	PermissionUpdate = "update"
	PermissionRead   = "read"
	PermissionDelete = "delete"
)

const (
	GroupRoleAdmin  = "admin"
	GroupRoleMember = "member"
)

const (
	GroupTypeAdmin = "admin" // admin, for /admin path with permissions check
	GroupTypeApp   = "app"
)

type Profile struct {
	Avatar  string         `json:"avatar,omitempty"`
	Gender  string         `json:"gender,omitempty"`
	City    string         `json:"city,omitempty"`
	Region  string         `json:"region,omitempty"`
	Country string         `json:"country,omitempty"`
	Extra   map[string]any `json:"extra,omitempty"`
}

func (p *Profile) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), p)
}

func (p Profile) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type User struct {
	ID        uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`

	Email       string     `json:"email" gorm:"size:128;uniqueIndex"`
	Password    string     `json:"-" gorm:"size:128"`
	Phone       string     `json:"phone,omitempty" gorm:"size:64;index"`
	FirstName   string     `json:"first_name,omitempty" gorm:"size:128"`
	LastName    string     `json:"last_name,omitempty" gorm:"size:128"`
	DisplayName string     `json:"display_name,omitempty" gorm:"size:128"`
	IsSuperUser bool       `json:"-"`
	IsStaff     bool       `json:"-"`
	Enabled     bool       `json:"-"`
	Activated   bool       `json:"-"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	LastLoginIP string     `json:"-" gorm:"size:128"`

	Source    string   `json:"-" gorm:"size:64;index"`
	Locale    string   `json:"locale,omitempty" gorm:"size:20"`
	Timezone  string   `json:"timezone,omitempty" gorm:"size:200"`
	Profile   *Profile `json:"profile,omitempty"`
	AuthToken string   `json:"token,omitempty" gorm:"-"`
}

// permission format
// users.read,users.create,users.update,users.delete, user.*
// pages.publish,pages.update,page.delete,page.*
type GroupPermission struct {
	Permissions []string
}

type Group struct {
	ID         uint            `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time       `json:"-" gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `json:"-"`
	Name       string          `json:"name" gorm:"size:200"`
	Type       string          `json:"type" gorm:"size:24;index"`
	Extra      string          `json:"extra"`
	Permission GroupPermission `json:"-"`
}

type GroupMember struct {
	ID        uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UserID    uint      `json:"-"`
	User      User      `json:"user"`
	GroupID   uint      `json:"-"`
	Group     Group     `json:"group"`
	Role      string    `json:"role" gorm:"size:60"`
}

func (u User) String() string {
	n := u.GetVisibleName()
	if n != "" {
		return fmt.Sprintf("%s(%s)", u.Email, n)
	}
	return u.Email
}

func (g Group) String() string {
	return fmt.Sprintf("%s(%d)", g.Name, g.ID)
}

func (u *User) GetVisibleName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	return u.LastName
}

func (u *User) GetProfile() Profile {
	if u.Profile != nil {
		return *u.Profile
	}
	return Profile{}
}

func (p *GroupPermission) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), p)
}

func (p GroupPermission) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// GenUniqueKey generate a unique value for a field in a table.
func GenUniqueKey(tx *gorm.DB, field string, size int) (key string) {
	key = random.RandText(size)
	for i := 0; i < 10; i++ {
		var c int64
		result := tx.Where(field, key).Limit(1).Count(&c)
		if result.Error != nil {
			break
		}
		if c > 0 {
			continue
		}
		return key
	}
	return ""
}
