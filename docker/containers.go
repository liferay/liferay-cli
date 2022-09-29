package docker

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
)

func BuildImages(verbose bool) {
	var s *spinner.Spinner

	if !verbose {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Building 'localdev' images..."
		s.FinalMSG = fmt.Sprintf("\u2705 Built 'localdev' images.\n")
		s.Start()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go BuildImage("dxp-server", path.Join(
		viper.GetString(constants.Const.RepoDir), "docker", "images", "dxp-server"),
		verbose, &wg)

	wg.Add(1)
	go BuildImage("localdev-server", path.Join(
		viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-server"),
		verbose, &wg)

	wg.Add(1)
	go BuildImage("localdev-dnsmasq", path.Join(
		viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-dnsmasq"),
		verbose, &wg)

	wg.Wait()

	if s != nil {
		s.Stop()
	}
}
