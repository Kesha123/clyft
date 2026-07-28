package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clyft",
	Short: "A GitOps tool for containers",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
