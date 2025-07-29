package post

import (
	"github.com/nolafw/config/pkg/config"
	"github.com/nolafw/config/pkg/runtimeconfig"
	postgrpc "github.com/nolafw/projecttemplate/internal/module/post/controller/grpc"
	"github.com/nolafw/projecttemplate/internal/module/post/controller/http"
	"github.com/nolafw/projecttemplate/internal/module/post/service"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	pb "github.com/nolafw/projecttemplate/service_adapter/post"
	"github.com/nolafw/rest/pkg/rest"
)

const ModuleName = "post"

func Params() (*runtimeconfig.Parameters, error) {
	return config.ModuleParams(ModuleName)
}

func NewModule(get *http.Get) *rest.Module {
	return &rest.Module{
		Path: "/post",
		Get:  get,
	}
}

func init() {
	dikit.AppendConstructors([]any{
		dikit.Bind[service.PostService](service.NewPostService),
		http.NewGet,
		dikit.AsModule(NewModule),
		// gRPC
		postgrpc.NewPostGRPCService,
		dikit.Bind[pb.PostServer](postgrpc.NewPostGRPCService),
	})
}
