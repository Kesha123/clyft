package cmd

import (
	"clyft/internal/artifact"
	"fmt"

	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Create local OCI layout artifact.",
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, _ := cmd.Flags().GetStringArray("tag")
		path, _ := cmd.Flags().GetStringSlice("path")
		if len(tag) == 0 {
			return fmt.Errorf("tag is required (pass -t/--tag)")
		}
		if len(path) == 0 {
			return fmt.Errorf("at least one path is required (pass -p/--path)")
		}

		ctx := cmd.Context()

		return artifact.Artifact(ctx, tag, path)
	},
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	artifactCmd.Flags().StringArrayP("tag", "t", []string{}, "OCI artifact tag (repeatable)")
	artifactCmd.Flags().StringSliceP("path", "p", []string{}, "Directory or file path.")
}
