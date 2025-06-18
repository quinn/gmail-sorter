/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/quinn/gmail-sorter/pkg/core"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("apply called")

		var err error
		api, err := gmailapi.Start()

		if err != nil {
			return err
		}

		db := db.NewDB()
		slog.Info("connected to database")

		defer func() {
			slog.Info("closing database")
			_ = db.Close()
		}()

		spec, err := core.NewSpec(api, db)

		if err != nil {
			return err
		}

		err = spec.Apply()

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
