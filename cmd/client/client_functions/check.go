package clientfunctions

import (
	"context"
	"new_filesync/config"
	"new_filesync/inretnal/fs"
	"new_filesync/proto"
	"path/filepath"
	"slices"
	"time"
)

func CheckFilesFromServer(cfg *config.Config, client_conn proto.SyncServiceClient) (*proto.FileListResponse, *proto.FileListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	serverFiles, err := client_conn.ListFiles(ctx, &proto.FileListRequest{})
	if err != nil {
		return nil, nil, err
	}

	indexOfGitignore := slices.IndexFunc(serverFiles.Files, func(file *proto.FileMeta) bool {
		return filepath.Base(file.Path) == ".gitignore"
	})

	serverFiles.Files = append(serverFiles.Files[:indexOfGitignore], serverFiles.Files[indexOfGitignore+1:]...)

	clientFiles, err := fs.ScanDir(cfg.MainPath)
	if err != nil {
		return nil, nil, err
	}

	filesToDownload := make([]*proto.FileMeta, 0, len(serverFiles.Files))
	filesToDelete := make([]*proto.FileMeta, 0, len(clientFiles.Files))

	for _, serverFile := range serverFiles.Files {
		if !slices.ContainsFunc(clientFiles.Files, func(file *proto.FileMeta) bool {
			return file.Hash == serverFile.Hash && file.Path == serverFile.Path
		}) {
			filesToDownload = append(filesToDownload, serverFile)
		}
	}

	for _, clientFile := range clientFiles.Files {
		if !slices.ContainsFunc(serverFiles.Files, func(file *proto.FileMeta) bool {
			return file.Hash == clientFile.Hash && file.Path == clientFile.Path
		}) {
			filesToDelete = append(filesToDelete, clientFile)
		}
	}

	return &proto.FileListResponse{Files: filesToDownload}, &proto.FileListResponse{Files: filesToDelete}, nil
}

func CheckFilesFromClient(cfg *config.Config, client_conn proto.SyncServiceClient) (*proto.FileListResponse, *proto.FileListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	serverFiles, err := client_conn.ListFiles(ctx, &proto.FileListRequest{})
	if err != nil {
		return nil, nil, err
	}

	indexOfGitignore := slices.IndexFunc(serverFiles.Files, func(file *proto.FileMeta) bool {
		return filepath.Base(file.Path) == ".gitignore"
	})
	serverFiles.Files = append(serverFiles.Files[:indexOfGitignore], serverFiles.Files[indexOfGitignore+1:]...)

	clientFiles, err := fs.ScanDir(cfg.MainPath)
	if err != nil {
		return nil, nil, err
	}

	filesToUpload := make([]*proto.FileMeta, 0, len(clientFiles.Files))
	filesToDelete := make([]*proto.FileMeta, 0, len(serverFiles.Files))

	for _, serverFile := range serverFiles.Files {
		if !slices.ContainsFunc(clientFiles.Files, func(file *proto.FileMeta) bool {
			return file.Hash == serverFile.Hash && file.Path == serverFile.Path
		}) {
			filesToDelete = append(filesToDelete, serverFile)
		}
	}

	for _, clientFile := range clientFiles.Files {
		if !slices.ContainsFunc(serverFiles.Files, func(file *proto.FileMeta) bool {
			return file.Hash == clientFile.Hash && file.Path == clientFile.Path
		}) {
			filesToUpload = append(filesToUpload, clientFile)
		}
	}

	return &proto.FileListResponse{Files: filesToUpload}, &proto.FileListResponse{Files: filesToDelete}, nil
}
