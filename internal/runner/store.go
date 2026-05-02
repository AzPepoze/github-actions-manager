package runner

import (
	"encoding/json"
	"os"
)

type Status string

const (
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusUnknown Status = "unknown"
)

type Runner struct {
	ProjectName string `json:"project_name"`
	ProjectURL  string `json:"project_url"`
	Token       string `json:"token"`
	InstallPath string `json:"install_path"`
	Status      Status `json:"status"`
}

type Store struct {
	path    string
	Runners []Runner `json:"runners"`
}

func NewStore(path string) (*Store, error) {
	s := &Store{
		path:    path,
		Runners: []Runner{},
	}
	if err := s.Load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, s)
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) Add(runner Runner) error {
	s.Runners = append(s.Runners, runner)
	return s.Save()
}
