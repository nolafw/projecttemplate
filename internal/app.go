package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/projecttemplate/internal/module/user"
	"github.com/nolafw/rest/pkg/mw"
	"github.com/nolafw/rest/pkg/pipeline"
	"github.com/nolafw/rest/pkg/rest"
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

	// TODO: ビルドしてdockerにmainだけを入れた場合、設定フォルダは全部参照できないので?
	//       docker build時は、フォルダ構成をそのままに、すべてのsettingフォルダのみコピーする。
	// TODO: module内と、query内のsettingフォルダを検索して、スライスに含める。
	paths := []string{
		"./internal",
	}
	// Run the app
	schema, params, err := config.Load(*env, "setting", paths)
	if err != nil {
		panic(err)
	}
	fmt.Printf("schema: %v\n", schema["default"])
	fmt.Printf("params: %v\n", params["default"])

	httpPipeline := CreateHttpPipeline()
	httpPipeline.Set()

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", params["default"].Server.Port),
	}
	server.ListenAndServe()
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
