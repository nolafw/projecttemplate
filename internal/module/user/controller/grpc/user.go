package grpc

import (
	"context"
	"fmt"

	"github.com/nolafw/projecttemplate/internal/module/user/service"
	pb "github.com/nolafw/projecttemplate/service_adapter/user"
)

// gRPCでの接続処理
type UserGRPCService struct {
	pb.UnimplementedUserServer
	service service.UserService
}

// DIで使う必要なserviceをinject
func NewUserGRPCService(service service.UserService) *UserGRPCService {
	return &UserGRPCService{
		service: service,
	}
}

// TODO:
// serviceは基本的に、こんな感じの実装をするルールにする
// context.Contextを引数に受け取る
// reqのところは、DTOで置き換える
// 返り値は、DTOとerror
func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	fmt.Println("GetUser called with request:", req.UserId)

	return &pb.GetUserResponse{
		UserId: "1",
		Name:   "John Doe",
		Email:  "j@example.com",
	}, nil
}
