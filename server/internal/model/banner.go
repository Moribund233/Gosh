package model

type Banner struct {
	BaseModel
	Title       string `gorm:"size:64" json:"title"`
	Subtitle    string `gorm:"size:128" json:"subtitle"`
	Description string `gorm:"size:256" json:"description"`
	Image       string `gorm:"size:256" json:"image"`
	Background  string `gorm:"size:128" json:"background"`
	Link        string `gorm:"size:256" json:"link"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Status      string `gorm:"size:20;default:on;index" json:"status"`
}

type BrandStory struct {
	BaseModel
	Title       string `gorm:"size:64;not null" json:"title"`
	Description string `gorm:"size:256" json:"description"`
	Logo        string `gorm:"size:256" json:"logo"`
	Badge       string `gorm:"size:32" json:"badge"`
	Link        string `gorm:"size:256" json:"link"`
	Status      string `gorm:"size:20;default:on" json:"status"`
}

const (
	StatusOn  = "on"
	StatusOff = "off"
)
