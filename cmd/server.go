package cmd

import (
	"os"

	"github.com/quinn/gmail-sorter/internal/web"
	"github.com/quinn/gmail-sorter/pkg/core"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Gmail API and DB as in cmd/apply.go
		api, err := gmailapi.Start()
		if err != nil {
			return err
		}

		db := db.NewDB()
		defer db.Close()
		spec, err := core.NewSpec(api, db)
		if err != nil {
			return err
		}

		server := web.NewServer(spec)
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
