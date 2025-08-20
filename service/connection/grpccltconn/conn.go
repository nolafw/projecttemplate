package grpccltconn

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nolafw/grpcgear/pkg/interceptor/metadata/metaclt"
	"github.com/nolafw/grpcgear/pkg/interceptor/retry/retryclt"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/projecttemplate/internal/plamo/logkit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: ここにサービスのコネクションを作成

// FIXME: サービスごとにコネクションが違う場合はどうするか?
func NewUserPostConnection() grpc.ClientConnInterface {
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
	host := "localhost:50051" // 環境変数から取得
	conn, err := grpc.NewClient(host, options...)

	// FIXME: errの時は、fxのinjectionはどうすればいいのか?
	if err != nil {
		panic(fmt.Sprintf("failed to connect to gRPC host %s: %v", host, err))
	}
	return conn
}

func init() {
	dikit.AppendConstructors([]any{
		NewUserPostConnection,
	})
}
