package goal

import (
	"errors"
	"time"

	"github.com/yourname/football-api/internal/domain"
)

var (
	ErrMatchNotFound      = errors.New("pertandingan tidak ditemukan")
	ErrMatchAlreadyDone   = errors.New("pertandingan sudah selesai, tidak bisa diubah")
	ErrPlayerNotInMatch   = errors.New("pemain bukan bagian dari tim yang bertanding")
	ErrScoreMismatch      = errors.New("jumlah gol tidak sesuai dengan skor yang dideklarasikan")
)

// MatchChecker dependency dari match module
type MatchChecker interface {
	GetMatchTeams(matchID uint) (homeTeamID, awayTeamID uint, matchDate time.Time, status domain.MatchStatus, err error)
	MarkCompleted(matchID uint, homeScore, awayScore int) error
}

// PlayerChecker dependency dari player module
type PlayerChecker interface {
	GetPlayerTeamID(playerID uint) (uint, error)
}

type Service interface {
	SubmitResult(req *SubmitResultRequest) ([]domain.Goal, error)
}

type service struct {
	repo          Repository
	matchChecker  MatchChecker
	playerChecker PlayerChecker
}

func NewService(repo Repository, matchChecker MatchChecker, playerChecker PlayerChecker) Service {
	return &service{repo, matchChecker, playerChecker}
}

func (s *service) SubmitResult(req *SubmitResultRequest) ([]domain.Goal, error) {
	// 1. Ambil info match
	homeTeamID, awayTeamID, _, status, err := s.matchChecker.GetMatchTeams(req.MatchID)
	if err != nil {
		return nil, err
	}
	if status == domain.MatchStatusCompleted {
		return nil, ErrMatchAlreadyDone
	}

	validTeams := map[uint]bool{homeTeamID: true, awayTeamID: true}

	// 2. Validasi setiap pencetak gol & hitung skor aktual
	homeGoalCount := 0
	awayGoalCount := 0
	goalsToCreate := make([]domain.Goal, 0, len(req.Goals))

	for _, gr := range req.Goals {
		playerTeamID, err := s.playerChecker.GetPlayerTeamID(gr.PlayerID)
		if err != nil {
			return nil, err
		}
		if !validTeams[playerTeamID] {
			return nil, ErrPlayerNotInMatch
		}

		if playerTeamID == homeTeamID {
			homeGoalCount++
		} else {
			awayGoalCount++
		}

		goalsToCreate = append(goalsToCreate, domain.Goal{
			MatchID:  req.MatchID,
			PlayerID: gr.PlayerID,
			TeamID:   playerTeamID,
			Minute:   gr.Minute,
		})
	}

	// 3. Validasi jumlah gol konsisten dengan skor
	if homeGoalCount != req.HomeScore || awayGoalCount != req.AwayScore {
		return nil, ErrScoreMismatch
	}

	// 4. Hapus gol lama kalau ada (untuk re-submit), lalu insert baru
	if err := s.repo.DeleteByMatchID(req.MatchID); err != nil {
		return nil, err
	}
	if err := s.repo.BulkCreate(goalsToCreate); err != nil {
		return nil, err
	}

	// 5. Update status match jadi completed + update skor
	if err := s.matchChecker.MarkCompleted(req.MatchID, req.HomeScore, req.AwayScore); err != nil {
		return nil, err
	}

	return s.repo.FindByMatchID(req.MatchID)
}
