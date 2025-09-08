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

		filesToDownload, filesToDelete, err := clientfunctions.CheckFilesFromServer(cfg, clientConn)
		if err != nil {
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
			log.Fatalln(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
