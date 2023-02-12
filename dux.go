package duxgo

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"os"
)

type Dux struct {
}

func New() *Dux {
	return &Dux{}
}

func (t Dux) Run() {

	var rootCmd = &cobra.Command{Use: "dux"}

	var subCmd = &cobra.Command{
		Use:   "sub [no options!]",
		Short: "My subcommand",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd PreRun with args: %v\n", args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd Run with args: %v\n", args)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd PostRun with args: %v\n", args)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside subCmd PersistentPostRun with args: %v\n", args)
		},
	}
	rootCmd.AddCommand(subCmd)

	var webCmd = &cobra.Command{
		Use:   "web",
		Short: "starting the web service",
		Run: func(cmd *cobra.Command, args []string) {
			color.Println(fmt.Sprintf("⇨ <red>%s</>", registry.Version))
		},
	}
	rootCmd.AddCommand(webCmd)

	var queueCmd = &cobra.Command{
		Use:   "queue",
		Short: "start queue service",
		Run: func(cmd *cobra.Command, args []string) {
			color.Println(fmt.Sprintf("⇨ <red>%s</>", registry.Version))
		},
	}
	rootCmd.AddCommand(queueCmd)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "View the version number",
		Run: func(cmd *cobra.Command, args []string) {
			color.Println(fmt.Sprintf("⇨ <red>%s</>", registry.Version))
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
