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

func NewGRPCServer() *grpc.Server {
	// TODO: interceptorを使って、リクエストのログを出力する
	// TODO: panicが起きたときの制御はどうなる?
	// そういった処理のセットを、httpPipelineのようにここの `opt`に渡すようにする

	return grpc.NewServer()
	// return nil // gRPCを使わない場合はnilを返す
}

func RegisterServerLifecycle(lc LC, srv *http.Server, grpcSrv *grpc.Server) *http.Server {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// HTTPサーバーを別goroutineで起動
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					// TODO: メッセージをわかりやすいものに変更する
					slog.Error("HTTP server ListenAndServe error", "error", err)
				}
			}()

			// gRPCサーバーを別goroutineで起動
			if grpcSrv != nil {
				go func() {
					// TODO: ポートを指定できるようにする
					listen, err := net.Listen("tcp", ":50051")
					if err != nil {
						// TODO: ログをちゃんとしたものに修正
						slog.Error("gRPC server failed to listen", "error", err)
						return
					}
					slog.Info("gRPC server starting on :50051")
					if err := grpcSrv.Serve(listen); err != nil {
						// TODO: ログをちゃんとしたものに修正
						slog.Error("gRPC server failed to start", "error", err)
					}
				}()
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if grpcSrv != nil {
				grpcSrv.Stop()
			}

			return srv.Shutdown(ctx)
		},
	})
	return srv
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
