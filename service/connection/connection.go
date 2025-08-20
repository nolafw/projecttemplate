package connection

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: ここにサービスのコネクションを作成

// FIXME: サービスごとにコネクションが違う場合はどうするか?
func NewUserPostConnection() *grpc.ClientConn {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	// FIXME: errの時はどうすればいいのか?
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

func init() {
	// TODO: ここでdiに注入
}
