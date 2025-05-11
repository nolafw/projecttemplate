package user

import (
	"github.com/nolafw/di/pkg/di"
	"github.com/nolafw/projecttemplate/internal/module/user/controller"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/rest/pkg/rest"
)

// TODO: nolacliでモジュールを作成したら、このファイルに
// 自動的に、NewModuleと、Constructorsを追加する
// さらに、module_provider.goにもConstructorsを追加すること

func NewModule(get *controller.Get, post *controller.Post) *rest.Module {
	return &rest.Module{
		Path: "/user",
		Get:  get,
		Post: post,
	}
}

func Constructors() []any {
	return []any{
		di.Bind[service.UserService](service.NewUserService),
		controller.NewGet,
		controller.NewPost,
		di.AsModule(NewModule),
	}
}
