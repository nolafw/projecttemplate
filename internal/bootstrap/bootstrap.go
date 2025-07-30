package bootstrap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/config/pkg/env"
	_ "github.com/nolafw/projecttemplate/internal/module"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/projecttemplate/internal/plamo/logkit"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: 別のファイルに分ける
type GlobalError struct {
	Message string `json:"message"`
}

// これを、cmd/main.goで実行する
func Run(envVal *string) {
	config.InitializeConfiguration(envVal, "./internal", "config")

	dikit.AppendConstructors([]any{
		NewHttpApp(envVal),
		NewGRPCApp(envVal),
		dikit.AsHttpPipeline(CreateHttpPipeline),
	})

	dikit.ProvideAndRun(dikit.Constructors(), dikit.RegisterGRPCServices(), false)
}

// HTTPサーバーの初期化
func NewHttpApp(envVal *string) func(lc dikit.LC, httpPipeline *pipeline.Http) *http.Server {
	return func(lc dikit.LC, httpPipeline *pipeline.Http) *http.Server {
		httpPipeline.Set()
		// TODO: envValを使うこと
		params, err := config.ModuleParams("default")
		if err != nil {
			log.Fatalf("default config parameters not found: %s", err)
		}
		// FIXME: 別の場所に移す
		logkit.SetLogLevel(params.Log.Level)

		httpSrv := &http.Server{
			Addr: fmt.Sprintf(":%d", params.Server.Port),
		}

		dikit.RegisterHTTPServerLifecycle(lc, httpSrv)
		return httpSrv
	}
}

func CreateHttpPipeline(modules []*rest.Module) *pipeline.Http {
	// TODO: 別のファイルに分ける
	panicResponse := &rest.Response{
		Code: http.StatusInternalServerError,
		Body: &GlobalError{Message: "internal server error"},
	}

	configParams, err := config.ModuleParams("default")
	if err != nil {
		log.Fatalf("default config parameters not found: %s", err)
	}
	cors := configParams.Cors

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

// gRPCサーバーの初期化
func NewGRPCApp(envVal *string) func(lc dikit.LC) *grpc.Server {
	return func(lc dikit.LC) *grpc.Server {
		// TODO: interceptorを使って、リクエストのログを出力する
		// TODO: panicが起きたときの制御はどうなる?
		// そういった処理のセットを、httpPipelineのようにここの `opt`に渡すようにする
		grpcSrv := grpc.NewServer()

		// reflectionは開発環境でのみ有効にする
		// TODO: config/env にIsLocal()を作って、それを使う
		if envVal != nil && (*envVal == string(env.Local) || *envVal == string(env.Develop)) {
			reflection.Register(grpcSrv)
			logkit.Info("gRPC reflection enabled for development environment", "env", *envVal)
		}

		dikit.RegisterGRPCServerLifecycle(lc, grpcSrv)
		return grpcSrv
	}
}

func logger(req *rest.Request, res *rest.Response) {
	logkit.Info(
		"TODO: メッセージ内容",
		"addr", req.RemoteAddr(),
		"method", req.Method(),
		"code", res.Code,
		"path", req.Path(),
		"user-agent", req.UserAgent(),
	)
}
