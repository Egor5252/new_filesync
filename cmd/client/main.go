package main

import (
	"context"
	"fmt"
	"log"
	"new_filesync/config"
	"new_filesync/inretnal/client"
	"new_filesync/proto"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%v", cfg.ServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client_conn := proto.NewSyncServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = filepath.WalkDir(cfg.MainPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // ошибка доступа к файлу
		}

		if d.IsDir() {
			return nil // пропускаем директории
		}

		if err := client.UploadFile(ctx, client_conn, cfg.MainPath, path); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}
}
