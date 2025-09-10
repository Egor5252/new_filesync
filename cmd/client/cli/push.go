package cli

import (
	"fmt"
	"log"
	clientfunctions "new_filesync/cmd/client/client_functions"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Синхронизация папки на сервере с локальной папкой",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Выполняется команда push")

		cfg, clientConn, clientClose, err := clientfunctions.StartClient()
		if err != nil {
			log.Println(err)
			return
		}
		defer clientClose()

		filesToUpload, filesToDelete, err := clientfunctions.CheckFilesFromClient(cfg, clientConn)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("Файлы для отправки")
		for _, file := range filesToUpload.Files {
			fmt.Println(file)
		}

		fmt.Println("Файлы для удаления")
		for _, file := range filesToDelete.Files {
			fmt.Println(file)
		}

		if err := clientfunctions.Push(cfg, clientConn); err != nil {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
