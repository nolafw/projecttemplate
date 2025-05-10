package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/projecttemplate/internal/module/user"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"
	"go.uber.org/fx"
)

type GlobalError struct {
	Message string `json:"message"`
}

var Modules = []rest.Module{
	user.Module,
}

func Register() {
	// ここで、module全体を合体させる。
}

// これを、cmd/main.goで実行する
func Run(env *string) {

	fx.New(
		fx.Provide(NewApp(env)),
		fx.Invoke(func(*http.Server) {}),
	).Run()

}

func NewApp(env *string) func(lc fx.Lifecycle) *http.Server {
	return func(lc fx.Lifecycle) *http.Server {
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

		httpPipeline := CreateHttpPipeline()
		httpPipeline.Set()
		srv := &http.Server{
			Addr: fmt.Sprintf(":%d", params["default"].Server.Port),
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				srv.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		})
		return srv
	}
}

func CreateHttpPipeline() *pipeline.Http {
	panicResponse := &rest.Response{
		Code:   http.StatusInternalServerError,
		Object: &GlobalError{Message: "internal server error"},
	}
	return &pipeline.Http{
		Modules: Modules,
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
