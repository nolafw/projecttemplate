package bootstrap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nolafw/config/pkg/config"
	_ "github.com/nolafw/projecttemplate/internal/module"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/projecttemplate/internal/plamo/logkit"
	pbPost "github.com/nolafw/projecttemplate/service_adapter/post"
	pbUser "github.com/nolafw/projecttemplate/service_adapter/user"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"

	"google.golang.org/grpc"
)

// TODO: 別のファイルに分ける
type GlobalError struct {
	Message string `json:"message"`
}

// これを、cmd/main.goで実行する
func Run(env *string) {
	config.InitializeConfiguration(env, "./internal", "config")

	dikit.AppendConstructors([]any{
		NewHttpApp(env),
		NewGRPCApp(env),
		dikit.AsHttpPipeline(CreateHttpPipeline),
	})

	dikit.ProvideAndRun(dikit.Constructors(), func(
		httpSrv *http.Server,
		grpcSrv *grpc.Server,
		// TODO: この引数を特定の型の配列にしてにしてinvokeしたい
		userAPI pbUser.UserServer,
		postAPI pbPost.PostServer,
	) {
		// TODO: ちゃんとgRPCが動いてるかチェック
		if grpcSrv != nil {
			// TODO: ここにgRPCのサービスを登録する
			pbUser.RegisterUserServer(grpcSrv, userAPI)
			pbPost.RegisterPostServer(grpcSrv, postAPI)
			fmt.Println("gRPC server registered!")
		}
	}, false)
}

// HTTPサーバーの初期化
func NewHttpApp(env *string) func(lc dikit.LC, httpPipeline *pipeline.Http) *http.Server {
	return func(lc dikit.LC, httpPipeline *pipeline.Http) *http.Server {
		httpPipeline.Set()
		// TODO: envを使うこと
		params, err := config.ModuleParams("default")
		if err != nil {
			log.Fatalf("default config parameters not found: %s", err)
		}

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
func NewGRPCApp(env *string) func(lc dikit.LC) *grpc.Server {
	return func(lc dikit.LC) *grpc.Server {
		// TODO: envを使うこと

		// TODO: interceptorを使って、リクエストのログを出力する
		// TODO: panicが起きたときの制御はどうなる?
		// そういった処理のセットを、httpPipelineのようにここの `opt`に渡すようにする
		grpcSrv := grpc.NewServer()
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
