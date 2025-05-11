package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/projecttemplate/internal/module/user"
	"github.com/nolafw/projecttemplate/internal/module/user/controller"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"
	"go.uber.org/fx"
)

type GlobalError struct {
	Message string `json:"message"`
}

func Register() {
	// ここで、module全体を合体させる。
}

// これを、cmd/main.goで実行する
func Run(env *string) {

	fx.New(
		fx.Provide(
			NewApp(env),
			fx.Annotate(service.NewUserService, fx.As(new(service.UserService))),
			controller.NewGet,
			controller.NewPost,
			AsModule(user.NewModule),
			fx.Annotate(CreateHttpPipeline, fx.ParamTags(`group:"modules"`)),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()

}

func AsModule(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"routes"`),
	)
}

func NewApp(env *string) func(lc fx.Lifecycle, httpPipeline *pipeline.Http) *http.Server {
	return func(lc fx.Lifecycle, httpPipeline *pipeline.Http) *http.Server {
		paths := []string{
			"./internal",
		}
		// Run the app
		schema, params, err := config.Load(*env, "config", paths)
		if err != nil {
			panic(err)
		}
		fmt.Printf("schema: %v\n", schema["default"])
		fmt.Printf("params: %v\n", params["default"])

		httpPipeline.Set()
		srv := &http.Server{
			Addr: fmt.Sprintf(":%d", params["default"].Server.Port),
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						// サーバー起動に失敗した場合のエラーログ
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
}

func CreateHttpPipeline(modules []*rest.Module) *pipeline.Http {
	panicResponse := &rest.Response{
		Code:   http.StatusInternalServerError,
		Object: &GlobalError{Message: "internal server error"},
	}
	return &pipeline.Http{
		Modules: modules,
		GlobalMiddlewares: []rest.Middleware{
			mw.VerifyBodyParsable,
		},
		PanicResponse: panicResponse,
		Logger:        CreateLogger(),
	}
}

func CreateLogger() func(req *rest.Request, res *rest.Response) {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(req *rest.Request, res *rest.Response) {
		// 出力先はファイルやlogstashに実装で変えれる。設定で変えれるようにしたほうがいいか?
		l.Info(
			"TODO: メッセージ内容",
			"addr", req.RemoteAddr(),
			"method", req.Method(),
			"code", res.Code,
			"path", req.Path(),
			"user-agent", req.UserAgent(),
		)
	}
}
