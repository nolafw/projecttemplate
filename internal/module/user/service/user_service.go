package service

// サービスを別のモジュールから使う場合は、
// 直接このサービスを呼び出すのではなく、
// 一度ServiceAdapterを通して呼び出すこと
// serviceの返す値は必ずDTOにすること
// entityを返さないように実装すること
// entityはserviceの中で処理でのみ使う。
type UserService interface {
}

func NewUserService() UserService {
	return &UserServiceImpl{}
}

type UserServiceImpl struct {
}
