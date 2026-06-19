package match

import (
	"errors"

	"github.com/yourname/football-api/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(match *domain.Match) error
	FindAll() ([]domain.Match, error)
	FindByID(id uint) (*domain.Match, error)
	Update(match *domain.Match) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(m *domain.Match) error {
	return r.db.Create(m).Error
}

func (r *repository) FindAll() ([]domain.Match, error) {
	var matches []domain.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").Find(&matches).Error
	return matches, err
}

func (r *repository) FindByID(id uint) (*domain.Match, error) {
	var m domain.Match
	err := r.db.Preload("HomeTeam").Preload("AwayTeam").
		Preload("Goals").Preload("Goals.Player").Preload("Goals.Team").
		First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (r *repository) Update(m *domain.Match) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&domain.Match{}, id).Error
}
