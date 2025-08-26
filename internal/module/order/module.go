package order

import (
	"github.com/nolafw/config/pkg/appconfig"
	"github.com/nolafw/config/pkg/registry"
	"github.com/nolafw/projecttemplate/internal/module/order/controller/http"
	"github.com/nolafw/projecttemplate/internal/module/order/service"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	"github.com/nolafw/rest/pkg/rest"
)

const ModuleName = "order"

func Params() (*appconfig.Parameters, error) {
	return registry.ModuleParams(ModuleName)
}

func NewModule(get *http.Get) *rest.Module {
	return &rest.Module{
		Path: "/order",
		Get:  get,
	}
}

func init() {
	dikit.AppendConstructors([]any{
		// FIXME: injectされてない
		dikit.Bind[service.OrderService](service.NewOrderService),
		http.NewGet,
		dikit.AsModule(NewModule),
	})
}
