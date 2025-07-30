package grpckit

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// panic情報を格納する構造体
type PanicInfo struct {
	Value      any
	Stack      string
	Method     string
	ClientInfo map[string]any
	Timestamp  time.Time
}

// スタックトレースを取得する関数
func getStackTrace() string {
	buf := make([]byte, 1024*4) // 4KB のバッファ
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// TODO: 何を出力するかは、クライアントコード側で決められるようにしたい。
// callback関数を引数に取って、そこで出力する内容はクライアントコードで決めるようにするか?
// panicリカバリ用のUnaryServerInterceptor
func RecoveryUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				// panic情報を収集
				panicInfo := PanicInfo{
					Value:      r,
					Stack:      getStackTrace(),
					Method:     info.FullMethod,
					ClientInfo: GetExtendedClientInfo(ctx),
					Timestamp:  time.Now(),
				}

				// panicログを出力
				logAttrs := []any{
					"panic_value", r,
					"method", info.FullMethod,
					"stack_trace", panicInfo.Stack,
					"timestamp", panicInfo.Timestamp,
				}

				// クライアント情報を追加
				for key, value := range panicInfo.ClientInfo {
					logAttrs = append(logAttrs, key, value)
				}

				logger.Error("gRPC handler panic recovered", logAttrs...)

				// 適切なgRPCエラーステータスを返す
				err = status.Error(codes.Internal, "internal server error")
				resp = nil
			}
		}()

		return handler(ctx, req)
	}
}

// TODO: 何を出力するかは、クライアントコード側で決められるようにしたい。
// callback関数を引数に取って、そこで出力する内容はクライアントコードで決めるようにするか?
// panicリカバリ用のStreamServerInterceptor
func RecoveryStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// panic情報を収集
				panicInfo := PanicInfo{
					Value:      r,
					Stack:      getStackTrace(),
					Method:     info.FullMethod,
					ClientInfo: GetExtendedClientInfo(stream.Context()),
					Timestamp:  time.Now(),
				}

				// panicログを出力
				logAttrs := []any{
					"panic_value", r,
					"method", info.FullMethod,
					"stack_trace", panicInfo.Stack,
					"timestamp", panicInfo.Timestamp,
					"client_stream", info.IsClientStream,
					"server_stream", info.IsServerStream,
				}

				// クライアント情報を追加
				for key, value := range panicInfo.ClientInfo {
					logAttrs = append(logAttrs, key, value)
				}

				logger.Error("gRPC stream handler panic recovered", logAttrs...)

				// 適切なgRPCエラーステータスを返す
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, stream)
	}
}
