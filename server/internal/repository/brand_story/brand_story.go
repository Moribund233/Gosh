package brand_story

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	Get() (*model.BrandStory, error)
	Upsert(story *model.BrandStory) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) Get() (*model.BrandStory, error) {
	var s model.BrandStory
	err := database.DB.Where("status = ?", model.StatusOn).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repo) Upsert(story *model.BrandStory) error {
	var existing model.BrandStory
	result := database.DB.Where("status = ?", model.StatusOn).First(&existing)
	if result.Error != nil {
		return database.DB.Create(story).Error
	}
	story.ID = existing.ID
	story.CreatedAt = existing.CreatedAt
	return database.DB.Save(story).Error
}
