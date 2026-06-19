package team

import (
	"errors"

	"github.com/yourname/football-api/internal/domain"
)

var (
	ErrTeamNotFound   = errors.New("tim tidak ditemukan")
	ErrDuplicateName  = errors.New("nama tim sudah digunakan")
)

type Service interface {
	Create(req *CreateTeamRequest) (*domain.Team, error)
	GetAll() ([]domain.Team, error)
	GetByID(id uint) (*domain.Team, error)
	Update(id uint, req *UpdateTeamRequest) (*domain.Team, error)
	Delete(id uint) error
	// Dipakai modul lain untuk validasi
	TeamExists(id uint) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(req *CreateTeamRequest) (*domain.Team, error) {
	exists, err := s.repo.ExistsByName(req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateName
	}

	team := &domain.Team{
		Name:                req.Name,
		Logo:                req.Logo,
		Founded:             req.Founded,
		HeadquartersAddress: req.HeadquartersAddress,
		HeadquartersCity:    req.HeadquartersCity,
	}
	if err := s.repo.Create(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *service) GetAll() ([]domain.Team, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id uint) (*domain.Team, error) {
	team, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}

func (s *service) Update(id uint, req *UpdateTeamRequest) (*domain.Team, error) {
	team, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		exists, err := s.repo.ExistsByName(req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDuplicateName
		}
		team.Name = req.Name
	}
	if req.Logo != "" {
		team.Logo = req.Logo
	}
	if req.Founded != 0 {
		team.Founded = req.Founded
	}
	if req.HeadquartersAddress != "" {
		team.HeadquartersAddress = req.HeadquartersAddress
	}
	if req.HeadquartersCity != "" {
		team.HeadquartersCity = req.HeadquartersCity
	}

	if err := s.repo.Update(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *service) TeamExists(id uint) (bool, error) {
	team, err := s.repo.FindByID(id)
	if err != nil {
		return false, err
	}
	return team != nil, nil
}
