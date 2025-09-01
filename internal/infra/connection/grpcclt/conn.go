package grpcclt

import (
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
)

func init() {
	dikit.AppendConstructors([]any{
		// gRPC clientは、型が`grpc.ClientConnInterface`で全て同じになってしまう
		// そのため、`dikit.ProvideNamed`を使って、名前を付けて区別する
		// クライアント側のコンストラクタでは、`dikit.InjectNamed`で、
		// インジェクトする名前を指定してそれぞれに適したコネクションを渡す
		dikit.ProvideNamed(NewUserPostConnection, PostConnName),
	})
}
