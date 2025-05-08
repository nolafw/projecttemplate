package internal

import (
	"fmt"
	"net/http"
	"time"

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

	panicResponse := &rest.Response{
		Code:   http.StatusInternalServerError,
		Object: &GlobalError{Message: "internal server error"},
	}
	httpPipeline := &pipeline.Http{
		Modules: Modules,
		GlobalMiddlewares: []rest.Middleware{
			mw.VerifyBodyParsable,
		},
		PanicResponse: panicResponse,
		Logger:        logIncomingRequest,
	}
	httpPipeline.Set()

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", params["default"].Server.Port),
	}
	server.ListenAndServe()
}

func logIncomingRequest(req *rest.Request, res *rest.Response) {
	currentTime := time.Now().Format("2006-01-02 15:04:05.000")
	l := fmt.Sprintf(
		"%s | %s | %s | %d | %s | %s",
		currentTime,
		req.RemoteAddr(),
		req.Method(),
		res.Code,
		req.Path(),
		req.UserAgent(),
	)
	// TODO: 出力先をファイルやlogstashに変えれるようにする
	fmt.Println(l)
}
