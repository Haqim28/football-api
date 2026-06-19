package match

import (
	"errors"
	"time"

	"github.com/yourname/football-api/internal/domain"
)

var (
	ErrMatchNotFound   = errors.New("pertandingan tidak ditemukan")
	ErrSameTeam        = errors.New("tim home dan away tidak boleh sama")
	ErrTeamNotFound    = errors.New("tim tidak ditemukan")
	ErrAlreadyCompleted = errors.New("pertandingan sudah selesai")
)

type TeamChecker interface {
	TeamExists(id uint) (bool, error)
}

type Service interface {
	Create(req *CreateMatchRequest) (*domain.Match, error)
	GetAll() ([]domain.Match, error)
	GetByID(id uint) (*domain.Match, error)
	Update(id uint, req *UpdateMatchRequest) (*domain.Match, error)
	Delete(id uint) error
	// Dipakai goal module
	GetMatchTeams(matchID uint) (homeTeamID, awayTeamID uint, matchDate time.Time, status domain.MatchStatus, err error)
	MarkCompleted(matchID uint, homeScore, awayScore int) error
}

type service struct {
	repo        Repository
	teamChecker TeamChecker
}

func NewService(repo Repository, teamChecker TeamChecker) Service {
	return &service{repo, teamChecker}
}

func (s *service) Create(req *CreateMatchRequest) (*domain.Match, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return nil, ErrSameTeam
	}

	for _, id := range []uint{req.HomeTeamID, req.AwayTeamID} {
		exists, err := s.teamChecker.TeamExists(id)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrTeamNotFound
		}
	}

	m := &domain.Match{
		MatchDate:  req.MatchDate,
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		Status:     domain.MatchStatusScheduled,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return s.repo.FindByID(m.ID)
}

func (s *service) GetAll() ([]domain.Match, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id uint) (*domain.Match, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMatchNotFound
	}
	return m, nil
}

func (s *service) Update(id uint, req *UpdateMatchRequest) (*domain.Match, error) {
	m, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if m.Status == domain.MatchStatusCompleted {
		return nil, ErrAlreadyCompleted
	}

	if req.HomeTeamID > 0 || req.AwayTeamID > 0 {
		newHome := m.HomeTeamID
		newAway := m.AwayTeamID
		if req.HomeTeamID > 0 {
			newHome = req.HomeTeamID
		}
		if req.AwayTeamID > 0 {
			newAway = req.AwayTeamID
		}
		if newHome == newAway {
			return nil, ErrSameTeam
		}
		for _, tid := range []uint{newHome, newAway} {
			exists, err := s.teamChecker.TeamExists(tid)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, ErrTeamNotFound
			}
		}
		m.HomeTeamID = newHome
		m.AwayTeamID = newAway
	}

	if req.MatchDate != nil {
		m.MatchDate = *req.MatchDate
	}

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}
	return s.repo.FindByID(m.ID)
}

func (s *service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *service) GetMatchTeams(matchID uint) (uint, uint, time.Time, domain.MatchStatus, error) {
	m, err := s.repo.FindByID(matchID)
	if err != nil {
		return 0, 0, time.Time{}, "", err
	}
	if m == nil {
		return 0, 0, time.Time{}, "", ErrMatchNotFound
	}
	return m.HomeTeamID, m.AwayTeamID, m.MatchDate, m.Status, nil
}

func (s *service) MarkCompleted(matchID uint, homeScore, awayScore int) error {
	m, err := s.repo.FindByID(matchID)
	if err != nil {
		return err
	}
	if m == nil {
		return ErrMatchNotFound
	}
	m.Status = domain.MatchStatusCompleted
	m.HomeScore = homeScore
	m.AwayScore = awayScore
	return s.repo.Update(m)
}
