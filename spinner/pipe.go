package spinner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/pkg/stdcopy"
)

func SpinnerPipe(s *spinner.Spinner, prefix string, verbose bool) func(io.ReadCloser) {
	return func(out io.ReadCloser) {
		if verbose {
			stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		} else if s != nil {
			c := make(chan (string))
			go func() {
				for {
					msg := <-c
					s.Suffix = fmt.Sprintf(prefix, msg)
				}
			}()

			reader := bufio.NewReader(out)
			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					close(c)
					break
				} else {
					c <- strings.TrimSpace(str)
				}
			}
		}
	}
}
