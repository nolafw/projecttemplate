package grpcclt

import (
	"context"

	"github.com/google/uuid"
	"github.com/nolafw/grpcgear/pkg/interceptor/metadata/metaclt"
	"github.com/nolafw/grpcgear/pkg/interceptor/retry/retryclt"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/projecttemplate/internal/plamo/logkit"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const PostConnName = "postConn"

func NewUserPostConnection(lc dikit.LC) (grpc.ClientConnInterface, error) {
	metadataConfig := &metaclt.ClientMetadataConfig{
		StaticMetadata: map[string]string{
			"x-api-version": "v1",
			"x-app-name":    "nolafw",
		},
		RequestIDKey: "x-request-id",
		RequestIdGenerator: func() string {
			return "req_" + uuid.New().String()
		},
		Logger: logkit.Logger(),
	}
	retryCofig := retryclt.DefaultConfig()

	options := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			metaclt.ClientMetadataUnaryInterceptor(metadataConfig),
			retryclt.UnaryClientInterceptor(retryCofig),
			// add other interceptors if needed
		),
		grpc.WithChainStreamInterceptor(
			metaclt.ClientMetadataStreamInterceptor(metadataConfig),
			retryclt.StreamClientInterceptor(retryCofig),
			// add other interceptors if needed
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	host := "localhost:50051" // TODO: 環境変数などから取得
	conn, err := grpc.NewClient(host, options...)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn, nil
}
