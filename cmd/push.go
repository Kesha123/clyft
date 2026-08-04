package cmd

import (
	"clyft/internal/push"
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local OCI layout to OCI registry.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("error: clyft push accepts exactly one argument.")
		}
		tag := args[0]
		ctx := cmd.Context()
		return push.Push(ctx, tag)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
