package docker

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"liferay.com/lcectl/ansicolor"
	"liferay.com/lcectl/constants"
)

func BuildImages(verbose bool) {
	var s *spinner.Spinner

	if !verbose {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Building 'localdev' images..."
		s.FinalMSG = fmt.Sprintf(ansicolor.Good + " 'localdev' images built.\n")
		s.Start()
		defer s.Stop()
	}

	var g errgroup.Group

	g.Go(func() error {
		return BuildImage("localdev-server", filepath.Join(
			viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-server"),
			verbose)
	})

	// g.Wait waits for all goroutines to complete
	// and returns the first non-nil error returned
	// by one of the goroutines.
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
