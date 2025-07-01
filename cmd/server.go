package cmd

import (
	"os"

	"github.com/quinn/gmail-sorter/internal/web"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := db.NewDB()
		if err != nil {
			return err
		}
		defer db.Close()

		server, err := web.NewServer(db)
		if err != nil {
			return err
		}

		err = server.Start(":" + os.Getenv("PORT"))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
