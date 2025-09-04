package clientfunctions

import (
	"context"
	"fmt"
	"new_filesync/config"
	"new_filesync/inretnal/client"
	"new_filesync/proto"
	"time"
)

func DownloadAll(cfg *config.Config, client_conn proto.SyncServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	files, err := client_conn.ListFiles(ctx, &proto.FileListRequest{})
	if err != nil {
		return err
	}

	for i, file := range files.Files {
		fmt.Printf("%v из %v\n", i+1, len(files.Files))
		if err := client.DownloadFile(ctx, client_conn, &proto.FileRequest{Path: file.Path}, cfg.MainPath); err != nil {
			return err
		}
	}

	return nil
}
