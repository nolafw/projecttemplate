package module

import (
	"github.com/nolafw/projecttemplate/internal/module/user"
)

func AllModules() [][]any {
	return [][]any{
		user.Constructors(),
	}
}
