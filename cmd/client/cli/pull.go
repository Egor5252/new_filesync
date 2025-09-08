package cli

import (
	"fmt"
	"log"
	clientfunctions "new_filesync/cmd/client/client_functions"
	"path/filepath"

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

		asd, err := clientfunctions.Check(cfg, clientConn)
		if err != nil {
			log.Fatalln(err)
		}

		for _, val := range asd.Files {
			fmt.Printf("Файл %v, путь %v, размер %v, последнее изменение %v, хеш %v\n", filepath.Base(val.Path), val.Path, val.Size, val.ModifiedUnix, val.Hash)
		}

		// if err := clientfunctions.DownloadAll(cfg, clientConn); err != nil {
		// 	log.Fatalln(err)
		// 	return
		// }
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
