package spinner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/pkg/stdcopy"
	"liferay.com/lcectl/ansicolor"
	"liferay.com/lcectl/docker"
)

func SpinnerPipe(s *spinner.Spinner, prefix string) func(io.ReadCloser, bool) {
	return func(out io.ReadCloser, verbose bool) {
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
					break
				} else {
					c <- ansicolor.StripCodes(
						strings.TrimSpace(
							string(
								docker.TrimLogHeader(bytes))))
				}
			}
		}
	}
}

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")]
}
