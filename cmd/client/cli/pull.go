package cli

import (
	"fmt"
	"log"
	clientfunctions "new_filesync/cmd/client/client_functions"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull files from server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Выполняется команда pull")

		cfg, clientConn, err := clientfunctions.StartClient()
		if err != nil {
			log.Fatalln(err)
			return
		}

		if err := clientfunctions.DownloadAll(cfg, clientConn); err != nil {
			log.Fatalln(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
