package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "opget",
	Short: color.New(color.FgCyan).Sprint("MadGet CLI"),
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}
