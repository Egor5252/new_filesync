package clientfunctions

import (
	"context"
	"new_filesync/config"
	"new_filesync/inretnal/fs"
	"new_filesync/proto"
	"slices"
	"time"
)

func Check(cfg *config.Config, client_conn proto.SyncServiceClient) (*proto.FileListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	serverFiles, err := client_conn.ListFiles(ctx, &proto.FileListRequest{})
	if err != nil {
		return nil, err
	}

	clientFiles, err := fs.ScanDir(cfg.MainPath)
	if err != nil {
		return nil, err
	}

	missingFiles := make([]*proto.FileMeta, 0)

	for _, serverFile := range serverFiles.Files {
		if !slices.ContainsFunc(clientFiles.Files, func(file *proto.FileMeta) bool {
			return file.Hash == serverFile.Hash && file.Path == serverFile.Path
		}) {
			missingFiles = append(missingFiles, serverFile)
		}
	}

	return &proto.FileListResponse{
		Files: missingFiles,
	}, nil
}
