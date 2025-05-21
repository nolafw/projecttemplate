package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nolafw/di/pkg/di"
	_ "github.com/nolafw/projecttemplate/internal/module"
	"github.com/nolafw/projecttemplate/internal/util"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"
	"go.uber.org/fx"
)

// TODO: 別のファイルに分ける
type GlobalError struct {
	Message string `json:"message"`
}

// これを、cmd/main.goで実行する
func Run(env *string) {
	util.LoadConfig(env)

	di.AppendConstructors([]any{
		NewApp(env),
		di.AsHttpPipeline(CreateHttpPipeline),
	})

	di.ProvideAndRun(di.Constructors(), func(*http.Server) {})
}

// lcを使って、http.Serverのライフサイクルをカスタマイズすることも可能
func NewApp(env *string) func(lc fx.Lifecycle, httpPipeline *pipeline.Http) *http.Server {
	return func(lc fx.Lifecycle, httpPipeline *pipeline.Http) *http.Server {

		httpPipeline.Set()
		srv := &http.Server{
			// TODO: paramsの値を渡す
			Addr: fmt.Sprintf(":%d", params["default"].Server.Port),
		}
		return di.RegisterHttpServerLifecycle(lc, srv)
	}
}

func CreateHttpPipeline(modules []*rest.Module) *pipeline.Http {
	// TODO: 別のファイルに分ける
	panicResponse := &rest.Response{
		Code:   http.StatusInternalServerError,
		Object: &GlobalError{Message: "internal server error"},
	}

	configParam, err := util.GetConfigParam("default")
	if err != nil {
		log.Fatalf("config not found: %s", err)
	}
	cors := configParam.Cors

	return &pipeline.Http{
		Modules: modules,
		GlobalMiddlewares: []rest.Middleware{
			mw.VerifyBodyParsable,
			mw.NewSimpleCors(cors),
		},
		PanicResponse: panicResponse,
		Logger:        logger,
	}
}

func logger(req *rest.Request, res *rest.Response) {
	util.Log().Info(
		"TODO: メッセージ内容",
		"addr", req.RemoteAddr(),
		"method", req.Method(),
		"code", res.Code,
		"path", req.Path(),
		"user-agent", req.UserAgent(),
	)
}
