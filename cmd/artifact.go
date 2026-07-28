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
		tag, _ := cmd.Flags().GetString("tag")
		path, _ := cmd.Flags().GetStringSlice("path")
		if tag == "" {
			return fmt.Errorf("tag is required (pass -t/--tag)")
		}
		if len(path) == 0 {
			return fmt.Errorf("at least one path is required (pass -p/--path)")
		}

		return artifact.Artifact(tag, path)
	},
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	artifactCmd.Flags().StringP("tag", "t", "", "OCI artifact tag")
	artifactCmd.Flags().StringSliceP("path", "p", []string{}, "Directory or file path.")
}
