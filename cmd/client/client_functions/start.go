package clientfunctions

import (
	"fmt"
	"new_filesync/config"
	"new_filesync/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartClient() (*config.Config, proto.SyncServiceClient, func() error, error) {
	cfg := config.MustLoad()

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%v", cfg.ServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, err
	}

	client_conn := proto.NewSyncServiceClient(conn)

	return cfg, client_conn, conn.Close, nil
}
