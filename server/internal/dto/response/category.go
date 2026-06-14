package response

import (
	"gosh/internal/model"
)

type CategoryResponse struct {
	ID        uint                `json:"id"`
	ParentID  *uint               `json:"parent_id"`
	Name      string              `json:"name"`
	Icon      string              `json:"icon"`
	Banner    string              `json:"banner"`
	SortOrder int                 `json:"sort_order"`
	Level     int                 `json:"level"`
	Children  []CategoryResponse  `json:"children,omitempty"`
}

func ToCategoryResponse(c *model.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		ParentID:  c.ParentID,
		Name:      c.Name,
		Icon:      c.Icon,
		Banner:    c.Banner,
		SortOrder: c.SortOrder,
		Level:     c.Level,
	}
}

func ToCategoryTree(categories []model.Category) []CategoryResponse {
	byParent := make(map[uint][]model.Category)
	for _, c := range categories {
		pid := uint(0)
		if c.ParentID != nil {
			pid = *c.ParentID
		}
		byParent[pid] = append(byParent[pid], c)
	}

	var build func(parentID uint, level int) []CategoryResponse
	build = func(parentID uint, level int) []CategoryResponse {
		var res []CategoryResponse
		for _, c := range byParent[parentID] {
			r := ToCategoryResponse(&c)
			r.Level = level
			r.Children = build(c.ID, level+1)
			res = append(res, r)
		}
		return res
	}

	return build(0, 0)
}
