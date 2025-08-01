package user

import (
	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/config/pkg/runtimeconfig"
	usergrpc "github.com/nolafw/projecttemplate/internal/module/user/controller/grpc"
	"github.com/nolafw/projecttemplate/internal/module/user/controller/http"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/rest/pkg/rest"

	pb "github.com/nolafw/projecttemplate/service_adapter/user"
)

// TODO: nolacliでモジュールを作成したら、このファイルに
// 自動的に、NewModuleと、Constructorsを追加する
// さらに、moduler.goにもimportを追加すること

const ModuleName = "user"

// TODO: 便利機能として、この関数も自動的にnolacliで生成する
func Params() (*runtimeconfig.Parameters, error) {
	return config.ModuleParams(ModuleName)
}

func NewModule(get *http.Get, post *http.Post) *rest.Module {
	return &rest.Module{
		Path: "/user",
		Get:  get,
		Post: post,
	}
}

func init() {
	dikit.AppendConstructors([]any{
		dikit.Bind[service.UserService](service.NewUserService),
		http.NewGet,
		http.NewPost,
		dikit.AsModule(NewModule),
		// gRPC
		dikit.AsGRPCService(usergrpc.NewUserGRPCService),
		dikit.Bind[pb.UserServer](usergrpc.NewUserGRPCService),
	})
}
