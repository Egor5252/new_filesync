package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"new_filesync/config"
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

	client := proto.NewSyncServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = filepath.WalkDir(cfg.MainPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // ошибка доступа к файлу
		}

		if d.IsDir() {
			return nil // пропускаем директории
		}

		stream, err := client.UploadFile(ctx)
		if err != nil {
			log.Println(err)
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		buffer := make([]byte, 32*1024) // 32 KB чанки — можно настроить

		for {
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				log.Fatalf("failed to read file: %v", err)
			}
			if n == 0 {
				break // EOF
			}

			// Отправляем чанк
			relPath, _ := filepath.Rel(cfg.MainPath, path)
			err = stream.Send(&proto.FileChunk{
				Path:    relPath,    // путь передаётся в каждом чанке (можно оптимизировать)
				Content: buffer[:n], // только считанная часть
			})
			if err != nil {
				log.Fatalf("failed to send chunk: %v", err)
			}
		}

		// Закрываем поток и получаем ответ от сервера
		resp, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Upload failed: %v", err)
		}

		log.Printf("Upload success: %v — %v", resp.Success, resp.Message)

		return nil
	})

	if err != nil {
		return
	}
}
