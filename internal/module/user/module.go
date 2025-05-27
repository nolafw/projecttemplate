package user

import (
	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/config/pkg/runtimeconfig"
	"github.com/nolafw/di/pkg/di"
	"github.com/nolafw/projecttemplate/internal/module/user/controller"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/rest/pkg/rest"
)

// TODO: nolacliでモジュールを作成したら、このファイルに
// 自動的に、NewModuleと、Constructorsを追加する
// さらに、moduler.goにもimportを追加すること

const ModuleName = "user"

// TODO: 便利機能として、この関数も自動的にnolacliで生成する
func Params() (*runtimeconfig.Parameters, error) {
	return config.Params(ModuleName)
}

func NewModule(get *controller.Get, post *controller.Post) *rest.Module {
	return &rest.Module{
		Path: "/user",
		Get:  get,
		Post: post,
	}
}

func init() {
	di.AppendConstructors([]any{
		di.Bind[service.UserService](service.NewUserService),
		controller.NewGet,
		controller.NewPost,
		di.AsModule(NewModule),
	})
}
