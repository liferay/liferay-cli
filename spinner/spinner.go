package spinner

import (
	"fmt"
	"io"
	"time"

	"github.com/briandowns/spinner"
)

type SpinOperation func(func(io.ReadCloser, bool)) int

func Spin(doing string, done string, verbose bool, operation SpinOperation) {
	var s *spinner.Spinner

	if !verbose {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = fmt.Sprintf(" %s 'localdev' environment...", doing)
		s.FinalMSG = fmt.Sprintf("\u2705 %s 'localdev' environment.\n", done)
		s.Start()
	}

	pipeSpinner := SpinnerPipe(s, fmt.Sprintf(" %s 'localdev' environment", done)+" [%s]")

	signal := operation(pipeSpinner)

	if s != nil {
		if signal > 0 {
			s.FinalMSG = fmt.Sprintf("\u2718 Something went wrong...\n")
		}

		s.Stop()
	}

}
