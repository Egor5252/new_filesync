package client

import (
	"context"
	"io"
	"log"
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

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("Upload success: %v — %v", resp.Success, resp.Message)

	return nil
}
