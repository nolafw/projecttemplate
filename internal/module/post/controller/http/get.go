package http

import (
	"net/http"

	"github.com/nolafw/projecttemplate/internal/module/post/service"
	"github.com/nolafw/rest/pkg/rest"
)

type Get struct {
	Service service.PostService
}

func NewGet(service service.PostService) *Get {
	return &Get{
		Service: service,
	}
}

func (c *Get) Handle(r *rest.Request) *rest.Response {
	return &rest.Response{
		Xml:        true,
		Code:       http.StatusOK,
		AddHeaders: map[string]string{"Server": "net/http"},
		Body:       c.Service.Anything(),
	}
}
