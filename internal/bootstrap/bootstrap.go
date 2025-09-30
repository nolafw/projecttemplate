package bootstrap

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/nolafw/config/pkg/env"
	"github.com/nolafw/config/pkg/registry"
	_ "github.com/nolafw/projecttemplate/internal/infra/connection/grpcclt"
	_ "github.com/nolafw/projecttemplate/internal/module"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/projecttemplate/internal/plamo/grpckit"
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
	registry.InitializeConfiguration(envVal, "./internal", "config")

	dikit.AppendConstructors([]any{
		NewHttpApp(envVal),
		NewGRPCApp(envVal),
		dikit.AsHttpPipeline(CreateHttpPipeline),
	})
	// TODO: putputFxLogは、環境変数で変えれるようにする
	dikit.ProvideAndRun(dikit.Constructors(), dikit.RegisterGRPCServices(), true)
}

// HTTPサーバーの初期化
func NewHttpApp(envVal *string) func(lc dikit.LC, httpPipeline *pipeline.Http) *http.Server {
	return func(lc dikit.LC, httpPipeline *pipeline.Http) *http.Server {
		httpPipeline.Set()
		// TODO: envValを使うこと
		params, err := registry.ModuleParams("default")
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

	configParams, err := registry.ModuleParams("default")
	if err != nil {
		log.Fatalf("default config parameters not found: %s", err)
	}
	cors := configParams.Cors

	return &pipeline.Http{
		Modules: modules,
		GlobalMiddlewares: []rest.Middleware{
			mw.Logging(logIncomingRequest),
			mw.RecoveryWithLogger(panicResponse, logPanicDetails),
			mw.VerifyBodyParsable,
			mw.NewSimpleCors(cors),
		},
	}
}

// gRPCサーバーの初期化
func NewGRPCApp(envVal *string) func(lc dikit.LC) *grpc.Server {
	return func(lc dikit.LC) *grpc.Server {
		// ログとpanicリカバリinterceptor付きのgRPCサーバーを作成
		grpcSrv := grpckit.NewGRPCServer(logkit.Logger())

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

func logIncomingRequest(req *rest.Request, res *rest.Response) {
	logkit.Info("HTTP Request",
		slog.String("method", req.Method()),
		slog.String("path", req.Path()),
		slog.Int("status_code", res.Code),
		slog.String("remote_addr", req.RemoteAddr()),
		slog.String("user_agent", req.UserAgent()),
		slog.String("type", "access_log"),
	)
}

func logPanicDetails(r *rest.Request, panicValue interface{}, stackTrace []byte) {
	logkit.Error("Panic Recovered",
		slog.String("method", r.Method()),
		slog.String("url", r.URLStr()),
		slog.String("remote_addr", r.RemoteAddr()),
		slog.String("user_agent", r.UserAgent()),
		slog.Any("panic_value", panicValue),
		slog.String("panic_type", fmt.Sprintf("%T", panicValue)),
		slog.String("stack_trace", string(stackTrace)),
		slog.String("type", "panic_log"),
	)
}
