package service

import (
	"context"
	"fmt"
	"io"
	"new_filesync/inretnal/fs"
	"new_filesync/proto"

	"google.golang.org/grpc"
)

type SyncServer struct {
	proto.UnimplementedSyncServiceServer
}

func (s *SyncServer) ListFiles(ctx context.Context, req *proto.FileListRequest) (*proto.FileListResponse, error) {
	files, err := fs.ScanDir("cmd/server/sync-data/source")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("ошибка сканирования основного каталога сервера")
	}

	return files, nil
}

func (s *SyncServer) UploadFile(stream grpc.ClientStreamingServer[proto.FileChunk, proto.UploadStatus]) error {
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
				return stream.SendAndClose(&proto.UploadStatus{
					Success: true,
					Message: "Файл записан и сохранён",
				})
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
			write, close, err = fs.FileWriter("cmd/server/sync-data/source", fileChunk.Path)
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
