package main

import (
	"fmt"
	"log"
	"net"
	"new_filesync/config"
	"new_filesync/inretnal/server"
	"new_filesync/proto"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.ServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterSyncServiceServer(s, &server.SyncServer{})
	log.Printf("Сервер запущен на порту :%v", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//==============================================================
