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
		NewApp(env),
		dikit.NewGRPCServer,
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

// lcを使って、http.Serverのライフサイクルをカスタマイズすることも可能
func NewApp(env *string) func(lc dikit.LC, httpPipeline *pipeline.Http, grpcSrv *grpc.Server) *http.Server {
	return func(lc dikit.LC, httpPipeline *pipeline.Http, grpcSrv *grpc.Server) *http.Server {

		httpPipeline.Set()

		params, err := config.ModuleParams("default")
		logkit.SetLogLevel(params.Log.Level)

		if err != nil {
			log.Fatalf("default config parameters not found: %s", err)
		}
		httpSrv := &http.Server{
			Addr: fmt.Sprintf(":%d", params.Server.Port),
		}

		// grpcが不要な場合は、nilを渡すことも可能。TODO: デフォルトではそうしておくこと。
		return dikit.RegisterServerLifecycle(lc, httpSrv, grpcSrv)
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
