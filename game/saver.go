package game

import (
	"os"
	"sort"
	"time"

	"github.com/henilmalaviya/filic"
	"github.com/henilmalaviya/golw/env"
	"github.com/henilmalaviya/golw/util"
)

type Snapshot struct {
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Stats     GameStats `json:"stats"`
	Grid      [][2]int  `json:"grid"`
}

/* -------------------------------------------------------------------------- */

type GameSaver struct {
	manager *Manager

	SaveDir      string
	SaveInterval time.Duration
	MaxSaves     int

	ticker *time.Ticker
}

func NewSaveManager(manager *Manager) *GameSaver {
	return &GameSaver{
		manager:      manager,
		SaveDir:      env.Get().SaveDirectory,
		SaveInterval: time.Second * time.Duration(env.Get().SaveInterval),
		MaxSaves:     env.Get().MaxSavesFiles,
		ticker:       nil,
	}
}

func (s *GameSaver) Snapshot() *Snapshot {
	game := s.manager.GetGame()

	return &Snapshot{
		Version:   1,
		Timestamp: time.Now(),
		Stats:     s.manager.GetStats(),
		Grid:      game.GetGrid().GetLiveCellCoordinates(),
	}
}

func (s *GameSaver) StartSaving() {
	if s.SaveInterval <= 0 {
		return // No saving needed
	}

	s.ticker = time.NewTicker(s.SaveInterval)
	go func() {
		for range s.ticker.C {

			snap := s.Snapshot()
			filePath := s.SaveDir + "/save_" + snap.Timestamp.Format("20060102_150405") + ".json"

			if err := util.WriteJSONToFile(filePath, snap); err != nil {
				util.GetLogger().Error("Failed to save game state", "error", err)
				return
			}
			util.GetLogger().Info("Game state saved", "file", filePath)

			if err := util.CleanupOldFiles(s.SaveDir, s.MaxSaves); err != nil {
				util.GetLogger().Error("Failed to cleanup old save files", "error", err)
			}
		}
	}()
}

func (s *GameSaver) StopSaving() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
}

func (s *GameSaver) getLatestSaveFile() (*filic.File, error) {
	dir := filic.NewDirectory(s.SaveDir)

	files, err := dir.ListFiles()
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil // No files found
	}

	// Sort files by modification time, newest first
	sort.Slice(files, func(i, j int) bool {
		infoI, errI := os.Stat(files[i].Path)
		infoJ, errJ := os.Stat(files[j].Path)
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	return files[0], nil
}

func (s *GameSaver) LoadLatestSnapshot() (*Snapshot, error) {
	latestFile, err := s.getLatestSaveFile()
	if err != nil {
		return nil, err
	}
	if latestFile == nil {
		return nil, nil // No save file found
	}

	var snap Snapshot
	if err := util.ReadJSONFromFile(latestFile.Path, &snap); err != nil {
		return nil, err
	}

	return &snap, nil
}

func (s *GameSaver) LoadSnapshot(snap *Snapshot) {
	if snap == nil {
		return
	}

	s.manager.mutex.Lock()
	defer s.manager.mutex.Unlock()

	s.manager.stats.LoadFromSnapshot(&snap.Stats)
	s.manager.game.GetGrid().Clear()
	for _, coord := range snap.Grid {
		s.manager.game.GetGrid().SetCell(coord[0], coord[1])
	}

	util.GetLogger().Info("Game state loaded from snapshot", "timestamp", snap.Timestamp, "generation", snap.Stats.Generation)
}

func (s *GameSaver) LoadLatest() error {
	snap, err := s.LoadLatestSnapshot()
	if err != nil {
		return err
	}
	if snap == nil {
		util.GetLogger().Info("No save file found to load")
		return nil
	}

	s.LoadSnapshot(snap)
	return nil
}
