package tenet

import (
	"fmt"

	"github.com/lingo-reviews/dev/api"
	tomb "gopkg.in/tomb.v1"
)

type dryRun struct{}

type dryRunService struct {
	tenetService
}

func (d *dryRun) Pull(bool) error {
	fmt.Println("Dry Run: Pulling tenets")
	return nil
}

func (d *dryRun) OpenService() (TenetService, error) {
	s := &dryRunService{}
	s.start()
	return s, nil
}

func (s *dryRunService) start() error {
	fmt.Println("Dry Run: Starting service...")
	return nil
}

func (s *dryRunService) Close() error {
	fmt.Println("Dry Run: Closing service")
	return nil
}

func (s *dryRunService) Review(filesc <-chan string, issuesc chan<- *api.Issue, t *tomb.Tomb) error {
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
