package main

import (
	"context"
	"fmt"
	"log"
	"new_filesync/config"
	"new_filesync/proto"
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

	files, err := client.ListFiles(ctx, &proto.FileListRequest{})
	if err != nil {
		log.Fatal("Ошибка получения списка файлов с сервера")
	}

	for _, file := range files.Files {
		fmt.Println(file)
	}
}
