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
	fmt.Println("Dry Run: Starting review...")

	for filename := range filesc {
		fmt.Printf("Dry Run: Reviewing: %s\n", filename)
		issuesc <- &api.Issue{
			Name:     "dryrun",
			Comment:  "Dry Run Issue",
			LineText: "Your code here",
			Position: &api.IssueRange{&api.Position{filename, 0, 1, 1}, &api.Position{filename, 0, 1, 1}},
		}
	}
	close(issuesc)

	fmt.Println("Dry Run: Finishing review")
	return nil
}

func (s *dryRunService) Info() (*api.Info, error) {
	return &api.Info{
		Name:        "dryrun",
		Usage:       "test lingo and configurations",
		Description: "test lingo and configurations ... description",
		Language:    "*",
		Version:     "0.1.0",
	}, nil
}

func (s *dryRunService) Language() (string, error) {
	return "*", nil
}
