package game

import (
	"sync"
	"time"

	"github.com/henilmalaviya/gol"
	"github.com/henilmalaviya/golw/util"
)

type GameStats struct {
	Generation int `json:"generation"`
	BirthCount int `json:"birth_count"`
	DeathCount int `json:"death_count"`

	mutex sync.Mutex `json:"-"`
}

func (s *GameStats) IncrementGeneration() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Generation++
}

func (s *GameStats) IncrementBirths(count int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.BirthCount += count
}

func (s *GameStats) IncrementDeaths(count int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.DeathCount += count
}

func (s *GameStats) Snapshot() GameStats {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return GameStats{
		Generation: s.Generation,
		BirthCount: s.BirthCount,
		DeathCount: s.DeathCount,
	}
}

func (s *GameStats) LoadFromSnapshot(snapshot *GameStats) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Generation = snapshot.Generation
	s.BirthCount = snapshot.BirthCount
	s.DeathCount = snapshot.DeathCount
}

/* -------------------------------------------------------------------------- */

type Manager struct {
	game  *gol.Game
	stats GameStats

	ticker *time.Ticker

	mutex sync.Mutex
}

func (m *Manager) GetGame() *gol.Game {
	return m.game
}

func (m *Manager) GetStats() GameStats {
	return m.stats.Snapshot()
}

func (m *Manager) Start(tickInterval time.Duration) {
	if m.ticker != nil {
		return // Already started
	}

	logger := util.GetLogger()
	logger.Info("Starting game tick loop", "interval_ms", tickInterval.Milliseconds())

	m.ticker = time.NewTicker(tickInterval)

	go func() {
		for range m.ticker.C {
			m.mutex.Lock()
			bornCells, diedCells := m.game.GetGrid().Tick()
			m.stats.IncrementGeneration()
			m.stats.IncrementBirths(len(bornCells))
			m.stats.IncrementDeaths(len(diedCells))
			m.mutex.Unlock()
		}
	}()
}

func (m *Manager) Stop() {
	if m.ticker == nil {
		return // Not started
	}

	logger := util.GetLogger()
	logger.Info("Stopping game tick loop")

	m.ticker.Stop()
	m.ticker = nil
}

func NewManager() *Manager {
	manager := &Manager{
		game:   gol.NewGame(),
		stats:  GameStats{},
		ticker: nil,
	}
	return manager
}
