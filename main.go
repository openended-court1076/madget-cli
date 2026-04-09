package main

import (
	"fmt"
	"os"

	"madget-cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	fmt.Println()
	root := cmd.GetRootCommand()
	root.CompletionOptions.DisableDefaultCmd = true
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.SetIn(os.Stdin)
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
