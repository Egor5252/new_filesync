package clientfunctions

import (
	"context"
	"fmt"
	"log"
	"new_filesync/config"
	"new_filesync/inretnal/client"
	"new_filesync/proto"
	"path/filepath"
	"time"
)

func Push(cfg *config.Config, client_conn proto.SyncServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	filesToUpload, filesToDelete, err := CheckFilesFromClient(cfg, client_conn)
	if err != nil {
		log.Println(err)
		return err
	}

	for i, file := range filesToUpload.Files {
		fmt.Printf("%v из %v\n", i+1, len(filesToUpload.Files))
		if err := client.UploadFile(ctx, client_conn, cfg.MainPath, filepath.Join(cfg.MainPath, file.Path)); err != nil {
			return err
		}
	}

	for i, file := range filesToDelete.Files {
		fmt.Printf("%v из %v\n", i+1, len(filesToDelete.Files))
		status, err := client_conn.DeleteFile(ctx, &proto.FileRequest{Path: file.Path})
		if err != nil {
			return err
		}
		fmt.Println(status)
	}

	return nil
}
