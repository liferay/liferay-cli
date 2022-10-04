package spinner

import (
	"fmt"
	"io"
	"time"

	"github.com/briandowns/spinner"
	"liferay.com/lcectl/ansicolor"
)

type SpinOptions struct {
	// The "doing" verb. e.g. "Building"
	Doing string
	// The "done" verb. e.g. "Built"
	Done string
	// The noun on which the operation is being performed. e.g. "rocket"
	On string
	// Whether to use spinner is enabled or not
	Enable bool
}

type SpinOperation func(func(io.ReadCloser, bool)) int

func Spin(options SpinOptions, operation SpinOperation) {
	var s *spinner.Spinner

	if options.Enable {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = fmt.Sprintf(" %s %s...", options.On, options.Doing)
		s.Start()
	}

	pipeSpinner := SpinnerPipe(s, fmt.Sprintf(" %s %s", options.On, options.Doing)+" [%s]")

	signal := operation(pipeSpinner)

	if s != nil {
		mark := ansicolor.Good
		msg := options.On + " " + options.Done + "."

		if s.FinalMSG != "" {
			if signal > 0 {
				mark = ansicolor.Bad
			}
			msg = s.FinalMSG
		}

		s.FinalMSG = fmt.Sprintf(mark+" %s\n", msg)

		s.Stop()
	}

}
