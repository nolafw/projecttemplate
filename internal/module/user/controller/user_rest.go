package controller

import (
	"fmt"
	"net/http"

	"github.com/nolafw/projecttemplate/internal/module/user/dto"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/projecttemplate/internal/module/user/vali"
	"github.com/nolafw/rest/pkg/rest"
	"github.com/nolafw/validator/pkg/rule"
	"github.com/nolafw/validator/pkg/validate"
)

// RESTでのresponse
// JSON or XML

type Get struct {
	Service *service.UserService
}

func NewGet(service *service.UserService) *Get {
	return &Get{
		Service: service,
	}
}

func (c *Get) Handle(r *rest.Request) *rest.Response {
	return &rest.Response{
		Xml:        true,
		Code:       http.StatusOK,
		AddHeaders: map[string]string{"Server": "net/http"},
		Object:     &dto.User{Id: 1, Name: "hoge"},
	}
}

type Post struct {
	Service *service.UserService
}

func NewPost(service *service.UserService) *Post {
	return &Post{
		Service: service,
	}
}

func (c *Post) Handle(r *rest.Request) *rest.Response {
	user, err := validate.HttpRequestBody[dto.User](r.Request(), &rule.RuleSet{
		Field: "Name",
		Rules: []rule.Predicate{
			vali.Required,
		},
	})
	if err != nil {
		// エラーレスポンスを返す
	}
	id, isEmpty := r.PathValue("id")
	if !isEmpty {
		// do something
		vE := validate.Map(
			map[string]any{"id": id},
			&rule.RuleSet{Field: "id", Rules: []rule.Predicate{vali.MaxLength}},
		)
		if vE != nil {
			// エラーレスポンスを返す
		}
	}

	fmt.Printf("user: %#v\n", user)
	// ここで、serviceを呼ぶ
	// user, err := c.Service.Create(user)

	return &rest.Response{
		Code:       http.StatusOK,
		AddHeaders: map[string]string{"Server": "net/http"},
		Object:     &dto.User{Id: 1, Name: "hoge"},
	}
}
