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
	"liferay.com/liferay/cli/ansicolor"
	"liferay.com/liferay/cli/docker"
	lstrings "liferay.com/liferay/cli/strings"
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
						s.Suffix = fmt.Sprintf(prefix, lstrings.TruncateText(msg, 80))
					}

					if exitPattern != "" {
						match, _ := regexp.MatchString(exitPattern, msg)
						if match {
							return -1
						}
					}
				}
			}
		}
		return 0
	}
}
