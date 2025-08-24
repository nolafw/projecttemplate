package order

type OrderServiceAdapter interface {
	GetOrder() (string, error)
}

// TODO: 将来的に、別アプリケーションに分けてRPCでのinjectする場合は
// どういう風になるか一度実験してみる
// type GetOrderRequest struct {
// 	OrderID string
// }

// type GetOrderResponse struct {
// 	OrderId string
// 	Amount  float64
// 	Status  string
// }

// type OrderServiceAdapter interface {
// 	GetOrder(req GetOrderRequest) (GetOrderResponse, error)
// }
