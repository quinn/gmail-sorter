package cmd

import (
	"github.com/quinn/gmail-sorter/internal/web"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"github.com/quinn/gmail-sorter/pkg/db"
	"github.com/quinn/gmail-sorter/pkg/core"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Gmail API and DB as in cmd/apply.go
		api, err := gmailapi.Start()
		if err != nil {
			panic(err)
		}
		db := db.NewDB()
		defer db.Close()
		spec, err := core.NewSpec(api, db)
		if err != nil {
			panic(err)
		}
		server := web.NewServer(spec)
		err = server.Start(":3000")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serverCmd.PersistentFlags().String("port", "", "the port to bind the server to")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
