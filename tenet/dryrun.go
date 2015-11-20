package tenet

import (
	"fmt"

	"github.com/lingo-reviews/dev/api"
)

type dryRun struct{}

type dryRunService struct{}

func (d *dryRun) Pull(bool) error {
	fmt.Println("Dry Run: Pulling tenets")
	return nil
}

func (d *dryRun) Service() (TenetService, error) {
	return &dryRunService{}, nil
}

func (s *dryRunService) Start() error {
	fmt.Println("Dry Run: Starting service...")
	return nil
}

func (s *dryRunService) Stop() error {
	fmt.Println("Dry Run: Stopping service")
	return nil
}

func (s *dryRunService) Review(filesc <-chan string, issuesc chan<- *api.Issue) error {
	return nil
}

func (s *dryRunService) Info() (*api.Info, error) {
	return &api.Info{Name: "dryrun"}, nil
}

func (s *dryRunService) Language() (string, error) {
	return "", nil
}
