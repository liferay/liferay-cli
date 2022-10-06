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
			c := make(chan (string))
			go func() {
				for {
					msg := <-c
					if msg != "" {
						s.FinalMSG = msg
						s.Suffix = fmt.Sprintf(prefix, truncateText(msg, 80))
					}
				}
			}()

			reader := bufio.NewReader(out)
			for {
				bytes, _, err := reader.ReadLine()

				if err != nil {
					close(c)
					return 1
				} else {
					msg := ansicolor.StripCodes(
						strings.TrimSpace(
							string(
								docker.TrimLogHeader(bytes))))

					c <- msg

					if !verbose && exitPattern != "" {
						match, _ := regexp.MatchString(exitPattern, msg)
						if match {
							return 300
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
