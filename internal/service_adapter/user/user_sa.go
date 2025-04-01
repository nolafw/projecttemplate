package user

import (
	"github.com/nolafw/projecttemplate/internal/module/user/service"
)

// service adapterのメソッドが返す値は必ずDTOにすること
// entityを返さないように実装すること
type UserSA interface{}

// 別のモジュールのserviceを、いったんこの
// ServiceAdapterを通して呼び出すことで、
// そのserviceの実装を隠蔽することができる
// こうすることで、最初はmodule/serviceをただラップして使っていたとして
// 後でマイクロサービスに移行した場合に、このserviceadapterの中で
// module/serviceを呼び出すのではなく、gRPCなどでマイクロサービスを呼び出すように変更するだけで済む
// ここは、モノリスからマイクロサービスに移行しやすくするための中間の層としての役割を持つ
// ServiceもServiceAdapterも、handlerにinjectして、そこで扱う
type UserSAImpl struct {
	userService service.UserService
}
