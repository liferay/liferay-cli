package spinner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/pkg/stdcopy"
	"liferay.com/lcectl/ansicolor"
	"liferay.com/lcectl/docker"
)

func SpinnerPipe(s *spinner.Spinner, prefix string) func(io.ReadCloser, bool, string) int {
	return func(out io.ReadCloser, verbose bool, exitPattern string) int {
		if verbose {
			stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		} else if s != nil {
			reader := bufio.NewReader(out)

			for {
				bytes, _, err := reader.ReadLine()

				if err == io.EOF {
					return 0
				} else if err != nil {
					return 1
				} else {
					msg := ansicolor.StripCodes(
						strings.TrimSpace(
							string(
								docker.TrimLogHeader(bytes))))

					if msg != "" {
						s.FinalMSG = msg
						s.Suffix = fmt.Sprintf(prefix, truncateText(msg, 80))
					}

					if exitPattern != "" {
						match, _ := regexp.MatchString(exitPattern, msg)
						if match {
							return 0
						}
					}
				}
			}
		}
		return 0
	}
}

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")]
}
