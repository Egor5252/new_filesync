package clientfunctions

import (
	"context"
	"fmt"
	"new_filesync/config"
	"new_filesync/inretnal/client"
	"new_filesync/proto"
	"os"
	"path/filepath"
	"time"
)

func Pull(cfg *config.Config, client_conn proto.SyncServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	filesToDownload, filesToDelete, err := CheckFilesFromServer(cfg, client_conn)
	if err != nil {
		return err
	}

	for i, file := range filesToDownload.Files {
		fmt.Printf("%v из %v\n", i+1, len(filesToDownload.Files))
		if err := client.DownloadFile(ctx, client_conn, &proto.FileRequest{Path: file.Path}, cfg.MainPath); err != nil {
			return err
		}
	}

	for i, file := range filesToDelete.Files {
		fmt.Printf("%v из %v\n", i+1, len(filesToDelete.Files))
		os.Remove(filepath.Join(cfg.MainPath, file.Path))
	}

	return nil
}
