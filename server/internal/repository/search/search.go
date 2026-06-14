package search

import (
	"gosh/internal/database"
	"gosh/internal/model"
)

type Repository interface {
	CreateHistory(history *model.SearchHistory) error
	FindHistoryByUserID(userID uint, limit int) ([]model.SearchHistory, error)
	DeleteHistoryByUserID(userID uint) error
	HotSearch(limit int) ([]model.HotSearch, error)
	IncrementHotSearch(query string) error
}

type repo struct{}

func New() Repository {
	return &repo{}
}

func (r *repo) CreateHistory(history *model.SearchHistory) error {
	return database.DB.Create(history).Error
}

func (r *repo) FindHistoryByUserID(userID uint, limit int) ([]model.SearchHistory, error) {
	var list []model.SearchHistory
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Find(&list).Error
	return list, err
}

func (r *repo) DeleteHistoryByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.SearchHistory{}).Error
}

func (r *repo) HotSearch(limit int) ([]model.HotSearch, error) {
	var list []model.HotSearch
	err := database.DB.Order("count desc").Limit(limit).Find(&list).Error
	return list, err
}

func (r *repo) IncrementHotSearch(query string) error {
	return database.DB.Exec("INSERT INTO hot_searches (query, count) VALUES (?, 1) ON CONFLICT(query) DO UPDATE SET count = count + 1", query).Error
}
