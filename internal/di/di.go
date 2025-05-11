package di

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/fx"
)

// TODO: 後でdiプロジェクトに移す

var constructors = []any{}

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

func ProvideAndRun(constructors []any, invocation any) {
	fx.New(
		fx.Provide(
			constructors...,
		),
		fx.Invoke(invocation),
	).Run()
}

func RegisterHttpServerLifecycle(lc fx.Lifecycle, srv *http.Server) *http.Server {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					// サーバー起動に失敗した場合のエラーログ
					// TODO: slogに変更する
					fmt.Printf("HTTP server ListenAndServe error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
