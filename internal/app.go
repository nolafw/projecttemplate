package internal

import (
	"fmt"

	"github.com/nolafw/config/pkg/config"
)

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
}
