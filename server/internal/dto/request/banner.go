package request

type CreateBannerRequest struct {
	Title       string `json:"title" binding:"omitempty,max=64"`
	Subtitle    string `json:"subtitle" binding:"omitempty,max=128"`
	Description string `json:"description" binding:"omitempty,max=256"`
	Image       string `json:"image" binding:"omitempty,max=256"`
	Background  string `json:"background" binding:"omitempty,max=128"`
	Link        string `json:"link" binding:"omitempty,max=256"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateBannerRequest struct {
	Title       string `json:"title" binding:"omitempty,max=64"`
	Subtitle    string `json:"subtitle" binding:"omitempty,max=128"`
	Description string `json:"description" binding:"omitempty,max=256"`
	Image       string `json:"image" binding:"omitempty,max=256"`
	Background  string `json:"background" binding:"omitempty,max=128"`
	Link        string `json:"link" binding:"omitempty,max=256"`
	SortOrder   *int   `json:"sort_order"`
	Status      string `json:"status" binding:"omitempty,oneof=on off"`
}

type UpdateBrandStoryRequest struct {
	Title       string `json:"title" binding:"omitempty,max=64"`
	Description string `json:"description" binding:"omitempty,max=256"`
	Logo        string `json:"logo" binding:"omitempty,max=256"`
	Badge       string `json:"badge" binding:"omitempty,max=32"`
	Link        string `json:"link" binding:"omitempty,max=256"`
}
