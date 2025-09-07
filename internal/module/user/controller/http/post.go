package http

import (
	"fmt"
	"net/http"

	"github.com/nolafw/projecttemplate/internal/module/user/dto"
	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/projecttemplate/internal/plamo/vkit"
	"github.com/nolafw/rest/pkg/rest"
	"github.com/nolafw/validator/pkg/rule"
)

type Post struct {
	Service service.UserService
}

func NewPost(service service.UserService) *Post {
	return &Post{
		Service: service,
	}
}

func (c *Post) Handle(r *rest.Request) *rest.Response {
	user, err := vkit.HttpRequestBody[dto.CreateUser](r, &rule.RuleSet{
		Field: "name",
		Rules: []rule.Rule{
			vkit.MaxLength(10),
		},
	})
	if err != nil {
		return &rest.Response{
			Code: http.StatusBadRequest,
			Body: err,
		}
	}
	id, isEmpty := r.PathValue("id")
	if !isEmpty {
		vE := vkit.Map(
			map[string]any{"id": id},
			&rule.RuleSet{Field: "id", Rules: []rule.Rule{vkit.MaxLength(10)}},
		)
		if vE != nil {
			return &rest.Response{
				Code: http.StatusBadRequest,
				Body: vE,
			}
		}
	}

	fmt.Printf("user: %#v\n", user)
	// ここで、serviceを呼ぶ
	// user, err := c.Service.Create(user)

	return &rest.Response{
		Code:       http.StatusOK,
		AddHeaders: map[string]string{"Server": "net/http"},
		Body:       &dto.CreateUser{Id: 1, Name: "hoge"},
	}
}
