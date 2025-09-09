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
		log.Println("Выполняется команда pull")

		cfg, clientConn, clientClose, err := clientfunctions.StartClient()
		if err != nil {
			log.Println(err)
			return
		}
		defer clientClose()

		filesToDownload, filesToDelete, err := clientfunctions.CheckFilesFromServer(cfg, clientConn)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("Файлы для загрузки")
		for _, file := range filesToDownload.Files {
			fmt.Println(file)
		}

		fmt.Println("Файлы для удаления")
		for _, file := range filesToDelete.Files {
			fmt.Println(file)
		}

		if err := clientfunctions.Pull(cfg, clientConn); err != nil {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
