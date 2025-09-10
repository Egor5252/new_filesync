package client

import (
	"context"
	"io"
	"new_filesync/inretnal/fs"
	"new_filesync/proto"
	"os"
	"path/filepath"
)

func UploadFile(ctx context.Context, client proto.SyncServiceClient, mainPath, fullPath string) error {
	stream, err := client.UploadFile(ctx)
	if err != nil {
		return err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 32*1024)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break // EOF
		}

		relPath, err := filepath.Rel(mainPath, fullPath)
		if err != nil {
			return err
		}

		if err := stream.Send(&proto.FileChunk{
			Path:    relPath,
			Content: buffer[:n],
		}); err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

// Usege
// err = filepath.WalkDir(cfg.MainPath, func(path string, d os.DirEntry, err error) error {
// 		if err != nil {
// 			return err // ошибка доступа к файлу
// 		}

// 		if d.IsDir() {
// 			return nil // пропускаем директории
// 		}

// 		if err := client.UploadFile(ctx, client_conn, cfg.MainPath, path); err != nil {
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return
// 	}

func DownloadFile(ctx context.Context, client_conn proto.SyncServiceClient, fileReq *proto.FileRequest, mainPath string) error {
	stream, err := client_conn.DownloadFile(ctx, fileReq)
	if err != nil {
		return err
	}

	var write func([]byte) error
	var close func() error

	for {
		fileChunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				if close != nil {
					err = close()
					if err != nil {
						return err
					}
				}
				return nil
			} else {
				if close != nil {
					err = close()
					if err != nil {
						return err
					}
				}
				return err
			}
		}

		if write == nil {
			write, close, err = fs.FileWriter(mainPath, filepath.Clean(fileChunk.Path))
			if err != nil {
				return err
			}
		}

		err = write(fileChunk.Content)
		if err != nil {
			if close != nil {
				err = close()
				if err != nil {
					return err
				}
			}
			return err
		}
	}
}
