package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clyft",
	Short: "A GitOps tool for containers",
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.Execute()
}
