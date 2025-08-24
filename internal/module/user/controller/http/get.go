package http

import (
	"net/http"

	"github.com/nolafw/projecttemplate/internal/module/user/dto"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/rest/pkg/rest"
)

// RESTでのresponse
// JSON or XML

type Get struct {
	Service service.UserService
}

func NewGet(service service.UserService) *Get {
	return &Get{
		Service: service,
	}
}

func (c *Get) Handle(r *rest.Request) *rest.Response {

	c.Service.Something() // DEBUG:

	return &rest.Response{
		Xml:        true,
		Code:       http.StatusOK,
		AddHeaders: map[string]string{"Server": "net/http"},
		Body:       &dto.User{Id: 1, Name: "hoge"},
	}
}
