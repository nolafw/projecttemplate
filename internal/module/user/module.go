package user

import (
	"github.com/nolafw/projecttemplate/internal/module/user/controller"
	"github.com/nolafw/rest/pkg/rest"
)

func NewModule(get *controller.Get, post *controller.Post) *rest.Module {
	return &rest.Module{
		Path: "/user",
		Get:  get,
		Post: post,
	}
}
