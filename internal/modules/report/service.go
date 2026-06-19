package report

import (
	"errors"
	"time"

	"github.com/yourname/football-api/internal/domain"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type MatchReportResponse struct {
	MatchID     uint       `json:"match_id"`
	MatchDate   time.Time  `json:"match_date"`
	HomeTeam    TeamInfo   `json:"home_team"`
	AwayTeam    TeamInfo   `json:"away_team"`
	HomeScore   int        `json:"home_score"`
	AwayScore   int        `json:"away_score"`
	Status      string     `json:"status"`
	MatchResult string     `json:"match_result"`
	TopScorer   *TopScorer `json:"top_scorer"`
	HomeWins    int        `json:"home_team_total_wins"`
	AwayWins    int        `json:"away_team_total_wins"`
	Goals       []GoalInfo `json:"goals"`
}

type TeamInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type TopScorer struct {
	PlayerID uint   `json:"player_id"`
	Name     string `json:"name"`
	TeamName string `json:"team_name"`
	Goals    int    `json:"goals"`
}

type GoalInfo struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	Minute     int    `json:"minute"`
}

// ─── Interfaces ───────────────────────────────────────────────────────────────

type MatchQuerier interface {
	FindAll() ([]domain.Match, error)
	FindByID(id uint) (*domain.Match, error)
}

// ─── Errors ───────────────────────────────────────────────────────────────────

var ErrMatchNotFound = errors.New("pertandingan tidak ditemukan")

// ─── Service ─────────────────────────────────────────────────────────────────

type Service interface {
	GetMatchReport(matchID uint) (*MatchReportResponse, error)
	GetAllReports() ([]MatchReportResponse, error)
}

type service struct {
	matchRepo MatchQuerier
}

func NewService(matchRepo MatchQuerier) Service {
	return &service{matchRepo}
}

func (s *service) GetMatchReport(matchID uint) (*MatchReportResponse, error) {
	m, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMatchNotFound
	}

	allMatches, err := s.matchRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return buildReport(m, allMatches), nil
}

func (s *service) GetAllReports() ([]MatchReportResponse, error) {
	allMatches, err := s.matchRepo.FindAll()
	if err != nil {
		return nil, err
	}

	reports := make([]MatchReportResponse, 0, len(allMatches))
	for i := range allMatches {
		r := buildReport(&allMatches[i], allMatches)
		reports = append(reports, *r)
	}
	return reports, nil
}

// ─── Builder ─────────────────────────────────────────────────────────────────

func buildReport(m *domain.Match, allMatches []domain.Match) *MatchReportResponse {
	report := &MatchReportResponse{
		MatchID:   m.ID,
		MatchDate: m.MatchDate,
		HomeScore: m.HomeScore,
		AwayScore: m.AwayScore,
		Status:    string(m.Status),
		Goals:     []GoalInfo{},
	}

	if m.HomeTeam != nil {
		report.HomeTeam = TeamInfo{ID: m.HomeTeam.ID, Name: m.HomeTeam.Name, Logo: m.HomeTeam.Logo}
	}
	if m.AwayTeam != nil {
		report.AwayTeam = TeamInfo{ID: m.AwayTeam.ID, Name: m.AwayTeam.Name, Logo: m.AwayTeam.Logo}
	}

	// Hasil pertandingan
	if m.Status == domain.MatchStatusCompleted {
		switch {
		case m.HomeScore > m.AwayScore:
			report.MatchResult = "Tim Home Menang"
		case m.AwayScore > m.HomeScore:
			report.MatchResult = "Tim Away Menang"
		default:
			report.MatchResult = "Draw"
		}
	} else {
		report.MatchResult = "Belum Selesai"
	}

	// Gol detail & top scorer
	scorerCount := make(map[uint]int)
	type scorerMeta struct{ name, team string }
	scorerMetas := make(map[uint]scorerMeta)

	for _, g := range m.Goals {
		gi := GoalInfo{Minute: g.Minute}
		if g.Player != nil {
			gi.PlayerID = g.Player.ID
			gi.PlayerName = g.Player.Name
		}
		if g.Team != nil {
			gi.TeamName = g.Team.Name
		}
		report.Goals = append(report.Goals, gi)
		scorerCount[gi.PlayerID]++
		if _, ok := scorerMetas[gi.PlayerID]; !ok {
			scorerMetas[gi.PlayerID] = scorerMeta{name: gi.PlayerName, team: gi.TeamName}
		}
	}

	// Top scorer
	maxGoals := 0
	var topScorerID uint
	for pid, cnt := range scorerCount {
		if cnt > maxGoals {
			maxGoals = cnt
			topScorerID = pid
		}
	}
	if maxGoals > 0 {
		meta := scorerMetas[topScorerID]
		report.TopScorer = &TopScorer{
			PlayerID: topScorerID,
			Name:     meta.name,
			TeamName: meta.team,
			Goals:    maxGoals,
		}
	}

	// Akumulasi kemenangan tim home & away dari semua match s.d. match ini
	homeWins := 0
	awayWins := 0
	for _, am := range allMatches {
		if am.Status != domain.MatchStatusCompleted || am.ID > m.ID {
			continue
		}
		// Kemenangan tim home (m.HomeTeamID) di semua match (bisa posisi home/away)
		if am.HomeTeamID == m.HomeTeamID && am.HomeScore > am.AwayScore {
			homeWins++
		} else if am.AwayTeamID == m.HomeTeamID && am.AwayScore > am.HomeScore {
			homeWins++
		}
		// Kemenangan tim away (m.AwayTeamID)
		if am.HomeTeamID == m.AwayTeamID && am.HomeScore > am.AwayScore {
			awayWins++
		} else if am.AwayTeamID == m.AwayTeamID && am.AwayScore > am.HomeScore {
			awayWins++
		}
	}
	report.HomeWins = homeWins
	report.AwayWins = awayWins

	return report
}
