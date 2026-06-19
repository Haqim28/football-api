package team

import (
	"errors"

	"github.com/yourname/football-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(team *domain.Team) error
	FindAll() ([]domain.Team, error)
	FindByID(id uint) (*domain.Team, error)
	Update(team *domain.Team) error
	Delete(id uint) error
	ExistsByName(name string, excludeID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(team *domain.Team) error {
	return r.db.Create(team).Error
}

func (r *repository) FindAll() ([]domain.Team, error) {
	var teams []domain.Team
	err := r.db.Find(&teams).Error
	return teams, err
}

func (r *repository) FindByID(id uint) (*domain.Team, error) {
	var team domain.Team
	err := r.db.First(&team, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &team, err
}

func (r *repository) Update(team *domain.Team) error {
	return r.db.Save(team).Error
}

func (r *repository) Delete(id uint) error {
	// Soft delete via GORM (set deleted_at)
	return r.db.Delete(&domain.Team{}, id).Error
}

func (r *repository) ExistsByName(name string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.Team{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}
