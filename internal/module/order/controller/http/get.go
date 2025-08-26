package http

import (
	"net/http"

	"github.com/nolafw/projecttemplate/internal/module/order/service"
	"github.com/nolafw/rest/pkg/rest"
)

type Get struct {
	Service service.OrderService
}

func NewGet(service service.OrderService) *Get {
	return &Get{
		Service: service,
	}
}

func (c *Get) Handle(r *rest.Request) *rest.Response {
	order, _ := c.Service.GetOrder() // DEBUG:

	return &rest.Response{
		Code: http.StatusOK,
		Body: order,
	}
}
