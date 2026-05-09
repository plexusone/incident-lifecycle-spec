// Package main provides the ilspec CLI for incident lifecycle management.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "ilspec",
		Short: "Incident lifecycle specification tools",
		Long:  "Tools for working with incident lifecycle artifacts (premortem, intra-mortem, postmortem)",
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(renderCmd())
	rootCmd.AddCommand(validateCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ilspec %s\n", version)
		},
	}
}
