package order

import "github.com/nolafw/projecttemplate/internal/module/order/dto"

type OrderService interface {
	// FIXME: 実際には引数や戻り値はオブジェクトになるので、そこをちゃんと実装して試してみる
	GetOrder() (*dto.Order, error)
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
