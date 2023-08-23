package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/constants"
	lstrings "liferay.com/liferay/cli/strings"
)

func BuildImages(verbose bool) {
	var s *spinner.Spinner

	image := viper.GetString(constants.Const.DockerLocaldevServerImage)

	if viper.GetBool(constants.Const.DockerLocaldevServerPullimage) {
		if !verbose {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Color("green")
			s.Suffix = " Pulling 'localdev' [" + image + "]..."
			s.FinalMSG = fmt.Sprintf(ansicolor.Good + " 'localdev' pulled [" + image + "]\n")
			s.Start()
			defer s.Stop()
		}

		dockerClient := GetDockerClient()

		readCloser, err := dockerClient.ImagePull(
			context.Background(), image, types.ImagePullOptions{})

		if err != nil {
			log.Fatal(err)
		}

		defer readCloser.Close()

		var out io.Writer = os.Stdout
		var fd = os.Stdout.Fd()
		isTerminal := true

		if !verbose {
			out = spinnerWriter{
				callback: func(bytes []byte) (int, error) {
					msg := ansicolor.StripCodes(
						strings.TrimSpace(
							string(TrimLogHeader(bytes))))

					if msg != "" {
						s.Suffix = fmt.Sprintf(
							" Pulling 'localdev' %s",
							lstrings.StripNewlines(lstrings.TruncateText(msg, 80)))
					}

					return 0, nil
				},
			}
			fd = 0
			isTerminal = true
		}

		err = jsonmessage.DisplayJSONMessagesStream(readCloser, out, fd, isTerminal, nil)

		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if !verbose {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Building 'localdev' image..."
		s.FinalMSG = fmt.Sprintf(ansicolor.Good + " 'localdev' image built.\n")
		s.Start()
		defer s.Stop()
	}

	var g errgroup.Group

	g.Go(func() error {
		return BuildImage(
			image,
			filepath.Join(
				viper.GetString(constants.Const.RepoDir), "docker", "images", "localdev-server"),
			verbose, s)
	})

	// g.Wait waits for all goroutines to complete
	// and returns the first non-nil error returned
	// by one of the goroutines.
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
