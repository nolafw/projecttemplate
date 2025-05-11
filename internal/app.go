package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/di/pkg/di"
	"github.com/nolafw/projecttemplate/internal/module"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"
	"go.uber.org/fx"
)

type GlobalError struct {
	Message string `json:"message"`
}

// これを、cmd/main.goで実行する
func Run(env *string) {

	for _, m := range module.AllModules() {
		di.AppendConstructors(m)
	}

	di.AppendConstructors([]any{
		NewApp(env),
		di.AsHttpPipeline(CreateHttpPipeline),
	})

	di.ProvideAndRun(di.Constructors(), func(*http.Server) {})
}

// lcを使って、http.Serverのライフサイクルをカスタマイズすることも可能
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
		// DEBUG:
		fmt.Printf("schema: %v\n", schema["default"])
		fmt.Printf("params: %v\n", params["default"])

		httpPipeline.Set()
		srv := &http.Server{
			Addr: fmt.Sprintf(":%d", params["default"].Server.Port),
		}
		return di.RegisterHttpServerLifecycle(lc, srv)
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
