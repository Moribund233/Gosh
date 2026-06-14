package model

type Category struct {
	BaseModel
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	Name      string    `gorm:"size:32;not null" json:"name"`
	Icon      string    `gorm:"size:32" json:"icon"`
	Banner    string    `gorm:"size:256" json:"banner"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Level     int       `gorm:"default:0" json:"level"`
	Children  []Category `gorm:"-" json:"children,omitempty"`
}
