package player

import (
	"errors"

	"github.com/yourname/football-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(player *domain.Player) error
	FindAll(teamID uint) ([]domain.Player, error)
	FindByID(id uint) (*domain.Player, error)
	Update(player *domain.Player) error
	Delete(id uint) error
	IsJerseyNumberTaken(teamID uint, jerseyNumber int, excludeID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(p *domain.Player) error {
	return r.db.Create(p).Error
}

func (r *repository) FindAll(teamID uint) ([]domain.Player, error) {
	var players []domain.Player
	q := r.db.Preload("Team")
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	err := q.Find(&players).Error
	return players, err
}

func (r *repository) FindByID(id uint) (*domain.Player, error) {
	var p domain.Player
	err := r.db.Preload("Team").First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *repository) Update(p *domain.Player) error {
	return r.db.Save(p).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&domain.Player{}, id).Error
}

// IsJerseyNumberTaken cek apakah nomor punggung sudah dipakai di tim yang sama
func (r *repository) IsJerseyNumberTaken(teamID uint, jerseyNumber int, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Player{}).
		Where("team_id = ? AND jersey_number = ?", teamID, jerseyNumber)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}
