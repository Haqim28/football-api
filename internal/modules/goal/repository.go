package goal

import (
	"github.com/yourname/football-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	BulkCreate(goals []domain.Goal) error
	FindByMatchID(matchID uint) ([]domain.Goal, error)
	DeleteByMatchID(matchID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) BulkCreate(goals []domain.Goal) error {
	if len(goals) == 0 {
		return nil
	}
	return r.db.Create(&goals).Error
}

func (r *repository) FindByMatchID(matchID uint) ([]domain.Goal, error) {
	var goals []domain.Goal
	err := r.db.Preload("Player").Preload("Team").
		Where("match_id = ?", matchID).
		Order("minute ASC").
		Find(&goals).Error
	return goals, err
}

func (r *repository) DeleteByMatchID(matchID uint) error {
	// Soft delete semua gol di match ini
	return r.db.Where("match_id = ?", matchID).Delete(&domain.Goal{}).Error
}
