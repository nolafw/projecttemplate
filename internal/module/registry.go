package module

import (
	"github.com/nolafw/projecttemplate/internal/module/user"
)

func Registry() [][]any {
	return [][]any{
		user.Constructors(),
	}
}
