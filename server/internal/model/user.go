package model

import "time"

type User struct {
	BaseModel
	Phone      string     `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Password   string     `gorm:"size:128;not null" json:"-"`
	Nickname   string     `gorm:"size:64" json:"nickname"`
	Avatar     string     `gorm:"size:256" json:"avatar"`
	Role       string     `gorm:"size:20;default:user;index" json:"role"`
	TenantID   *uint      `gorm:"index" json:"tenant_id"`
	Status     string     `gorm:"size:20;default:active" json:"status"`
	Points     int        `gorm:"default:0" json:"points"`
	LastLogin  *time.Time `json:"last_login"`
}

const (
	RoleUser       = "user"
	RoleMerchant   = "merchant"
	RoleSupport    = "support"
	RoleOperator   = "operator"
	RoleSuperAdmin = "super_admin"

	StatusActive   = "active"
	StatusDisabled = "disabled"
)
