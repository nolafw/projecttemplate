package user

import (
	"github.com/nolafw/projecttemplate/internal/module/user/controller"
	"github.com/nolafw/rest/pkg/rest"
)

var Module = rest.Module{
	Path: "/user",
	Get:  &controller.Get{},
	Post: &controller.Post{},
}
