package duxgo

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type Dux struct {
}

func New() *Dux {
	return &Dux{}
}

func (t Dux) run() {

	var cmd = &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
