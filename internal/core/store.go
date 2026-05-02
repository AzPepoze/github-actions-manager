package core

type Store struct {
	Runners []Runner
}

func NewStore() *Store {
	return &Store{
		Runners: []Runner{},
	}
}

func (s *Store) SetRunners(runners []Runner) {
	s.Runners = runners
}

func (s *Store) UpdateStatus(index int, status Status) error {
	if index < 0 || index >= len(s.Runners) {
		return nil
	}
	s.Runners[index].Status = status
	return nil
}

func (s *Store) UpdateInstalled(index int, installed bool) error {
	if index < 0 || index >= len(s.Runners) {
		return nil
	}
	s.Runners[index].IsInstalled = installed
	return nil
}
