/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/quinn/gmail-sorter/internal/tui"
	"github.com/quinn/gmail-sorter/pkg/handlers"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Run the TUI",
	Long:  "Run the TUI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tui called")
		// Start at the menu page
		startLink := handlers.IndexAction.Link()

		fmt.Println("startLink: ", startLink.Action().ID)
		// Create Bubble Tea program with initial page
		prog := tui.NewProgram(tui.Page{Current: startLink})

		// Run program (blocks until quit)
		if _, err := prog.Run(); err != nil {
			log.Fatalf("bubbletea run: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
