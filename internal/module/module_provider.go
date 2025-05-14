package module

import (
	"github.com/nolafw/projecttemplate/internal/module/user"
)

// FIXME: これはやめて、di側に移したい。
// 各moduleのファイル内で
// di.RegisterConstructors([]any{})にしたい。
// そこで登録したものを、app.goで呼び出したい
func AllModules() [][]any {
	return [][]any{
		user.Constructors(),
	}
}
