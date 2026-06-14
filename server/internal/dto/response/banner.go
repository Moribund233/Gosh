package response

import (
	"gosh/internal/model"
)

type BannerResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Background  string `json:"background"`
	Link        string `json:"link"`
	SortOrder   int    `json:"sort_order"`
	Status      string `json:"status"`
}

type BrandStoryResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Badge       string `json:"badge"`
	Link        string `json:"link"`
}

func ToBannerResponse(b *model.Banner) BannerResponse {
	return BannerResponse{
		ID:          b.ID,
		Title:       b.Title,
		Subtitle:    b.Subtitle,
		Description: b.Description,
		Image:       b.Image,
		Background:  b.Background,
		Link:        b.Link,
		SortOrder:   b.SortOrder,
		Status:      b.Status,
	}
}

func ToBannerList(banners []model.Banner) []BannerResponse {
	list := make([]BannerResponse, len(banners))
	for i, b := range banners {
		list[i] = ToBannerResponse(&b)
	}
	return list
}

func ToBrandStoryResponse(s *model.BrandStory) BrandStoryResponse {
	return BrandStoryResponse{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		Logo:        s.Logo,
		Badge:       s.Badge,
		Link:        s.Link,
	}
}
