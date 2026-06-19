package player

import (
	"errors"

	"github.com/yourname/football-api/internal/domain"
)

var (
	ErrPlayerNotFound    = errors.New("pemain tidak ditemukan")
	ErrDuplicateJersey   = errors.New("nomor punggung sudah digunakan di tim ini")
	ErrInvalidPosition   = errors.New("posisi tidak valid, harus: penyerang, gelandang, bertahan, atau penjaga_gawang")
	ErrTeamNotFound      = errors.New("tim tidak ditemukan")
)

var validPositions = map[domain.PlayerPosition]bool{
	domain.PositionForward:    true,
	domain.PositionMidfielder: true,
	domain.PositionDefender:   true,
	domain.PositionGoalkeeper: true,
}

// TeamChecker dependency dari team module
type TeamChecker interface {
	TeamExists(id uint) (bool, error)
}

type Service interface {
	Create(req *CreatePlayerRequest) (*domain.Player, error)
	GetAll(teamID uint) ([]domain.Player, error)
	GetByID(id uint) (*domain.Player, error)
	Update(id uint, req *UpdatePlayerRequest) (*domain.Player, error)
	Delete(id uint) error
	// Dipakai goal module
	GetPlayerTeamID(playerID uint) (uint, error)
}

type service struct {
	repo        Repository
	teamChecker TeamChecker
}

func NewService(repo Repository, teamChecker TeamChecker) Service {
	return &service{repo, teamChecker}
}

func (s *service) Create(req *CreatePlayerRequest) (*domain.Player, error) {
	if !validPositions[req.Position] {
		return nil, ErrInvalidPosition
	}

	exists, err := s.teamChecker.TeamExists(req.TeamID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTeamNotFound
	}

	taken, err := s.repo.IsJerseyNumberTaken(req.TeamID, req.JerseyNumber, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrDuplicateJersey
	}

	p := &domain.Player{
		TeamID:       req.TeamID,
		Name:         req.Name,
		Height:       req.Height,
		Weight:       req.Weight,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetAll(teamID uint) ([]domain.Player, error) {
	return s.repo.FindAll(teamID)
}

func (s *service) GetByID(id uint) (*domain.Player, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrPlayerNotFound
	}
	return p, nil
}

func (s *service) Update(id uint, req *UpdatePlayerRequest) (*domain.Player, error) {
	p, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Position != "" && !validPositions[req.Position] {
		return nil, ErrInvalidPosition
	}

	// Kalau pindah tim, validasi tim tujuan
	targetTeamID := p.TeamID
	if req.TeamID > 0 && req.TeamID != p.TeamID {
		exists, err := s.teamChecker.TeamExists(req.TeamID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrTeamNotFound
		}
		targetTeamID = req.TeamID
		p.TeamID = req.TeamID
	}

	// Kalau jersey berubah, cek duplikasi di tim target
	newJersey := p.JerseyNumber
	if req.JerseyNumber > 0 && req.JerseyNumber != p.JerseyNumber {
		newJersey = req.JerseyNumber
	}
	taken, err := s.repo.IsJerseyNumberTaken(targetTeamID, newJersey, id)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrDuplicateJersey
	}

	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Height > 0 {
		p.Height = req.Height
	}
	if req.Weight > 0 {
		p.Weight = req.Weight
	}
	if req.Position != "" {
		p.Position = req.Position
	}
	if req.JerseyNumber > 0 {
		p.JerseyNumber = req.JerseyNumber
	}

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *service) GetPlayerTeamID(playerID uint) (uint, error) {
	p, err := s.repo.FindByID(playerID)
	if err != nil {
		return 0, err
	}
	if p == nil {
		return 0, ErrPlayerNotFound
	}
	return p.TeamID, nil
}
