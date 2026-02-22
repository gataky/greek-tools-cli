package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "greekmaster",
	Short: "Greek Case Master - Learn Modern Greek noun declension",
	Long: `Greek Case Master is a CLI educational tool that helps English speakers
learning Modern Greek transition from rote memorization of noun declension
tables to instinctive application of cases in real sentence contexts.`,
	Version: version,
}

func init() {
	// Add commands will be registered here as we implement them
	// rootCmd.AddCommand(importCmd)
	// rootCmd.AddCommand(practiceCmd)
	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
