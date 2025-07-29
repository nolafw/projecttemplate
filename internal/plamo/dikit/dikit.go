package dikit

import (
	"context"
	"net"
	"net/http"

	"log/slog"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var constructors = []any{}

type LC = fx.Lifecycle

func AppendConstructors(adding []any) error {
	constructors = append(constructors, adding...)
	return nil
}

func Constructors() []any {
	return constructors
}

func AsModule(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"modules"`),
	)
}

func AsHttpPipeline(f any) any {
	return fx.Annotate(f, fx.ParamTags(`group:"modules"`))
}

func Bind[T any](concrete any) any {
	return fx.Annotate(concrete, fx.As(new(T)))
}

func ProvideAndRun(constructors []any, invocation any, outputFxLog bool) {
	options := []fx.Option{
		fx.Provide(
			constructors...,
		),
		fx.Invoke(invocation),
	}

	if !outputFxLog {
		options = append(options, fx.NopLogger)
	}

	fx.New(options...).Run()
}

func RegisterHTTPServerLifecycle(lc LC, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				slog.Info("HTTP server starting", "addr", srv.Addr)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("HTTP server ListenAndServe error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down HTTP server")
			return srv.Shutdown(ctx)
		},
	})
}

func RegisterGRPCServerLifecycle(lc LC, grpcSrv *grpc.Server) {
	if grpcSrv == nil {
		return
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				// TODO: ポートを設定可能にする
				listen, err := net.Listen("tcp", ":50051")
				if err != nil {
					slog.Error("gRPC server failed to listen", "error", err)
					return
				}
				slog.Info("gRPC server starting", "addr", ":50051")
				if err := grpcSrv.Serve(listen); err != nil {
					slog.Error("gRPC server failed to start", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down gRPC server")
			grpcSrv.GracefulStop()
			return nil
		},
	})
}

// | 用途           | ライブラリ                                                              |
// | ------------ | ------------------------------------------------------------------ |
// | 認証・認可        | [grpc\_auth](https://github.com/grpc-ecosystem/go-grpc-middleware) |
// | ログ出力         | `zap` + `grpc_zap`                                                 |
// | リカバリ（panic）  | `grpc_recovery`                                                    |
// | リクエストバリデーション | `protoc-gen-validate`                                              |

// ✅ 自作Interceptorの活用パターン
// 以下のような処理が共通化可能です：
//     認証（JWTトークンの検証）
//     メタデータ（ヘッダー）の検査・追加
//     ログ出力
//     メトリクス収集（Prometheusなど）
//     トレース（OpenTelemetry）
