package grpckit

import (
	"log/slog"

	"google.golang.org/grpc"
)

// gRPCサーバー作成用のヘルパー関数
func NewGRPCServerWithLogging(logger *slog.Logger, opts ...grpc.ServerOption) *grpc.Server {
	// ログ出力interceptorを追加
	defaultOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(LoggingUnaryInterceptor(logger)),
		grpc.StreamInterceptor(LoggingStreamInterceptor(logger)),
	}

	// 追加のオプションをマージ
	allOpts := append(defaultOpts, opts...)

	return grpc.NewServer(allOpts...)
}

// panicリカバリ機能付きgRPCサーバー作成用のヘルパー関数
func NewGRPCServerWithRecovery(logger *slog.Logger, opts ...grpc.ServerOption) *grpc.Server {
	// panicリカバリinterceptorを追加
	defaultOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(RecoveryUnaryInterceptor(logger)),
		grpc.StreamInterceptor(RecoveryStreamInterceptor(logger)),
	}

	// 追加のオプションをマージ
	allOpts := append(defaultOpts, opts...)

	return grpc.NewServer(allOpts...)
}

// ログ出力とpanicリカバリ両方を含むgRPCサーバー作成用のヘルパー関数
func NewGRPCServerWithLoggingAndRecovery(logger *slog.Logger, opts ...grpc.ServerOption) *grpc.Server {
	// チェーン化された複数のinterceptorを作成
	// 注意: interceptorの順序は重要です。recoveryは最初に（最外側に）配置すべきです
	defaultOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			RecoveryUnaryInterceptor(logger), // 最外側: panic を最初にキャッチ
			LoggingUnaryInterceptor(logger),  // 内側: ログ出力
		),
		grpc.ChainStreamInterceptor(
			RecoveryStreamInterceptor(logger), // 最外側: panic を最初にキャッチ
			LoggingStreamInterceptor(logger),  // 内側: ログ出力
		),
	}

	// 追加のオプションをマージ
	allOpts := append(defaultOpts, opts...)

	return grpc.NewServer(allOpts...)
}

// | 用途           | ライブラリ                                                              |
// | ------------ | ------------------------------------------------------------------ |
// | 認証・認可        | [grpc\_auth](https://github.com/grpc-ecosystem/go-grpc-middleware) |

// 認証（JWTトークンの検証）
// メタデータ（ヘッダー）の検査・追加
