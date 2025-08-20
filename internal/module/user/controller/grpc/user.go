package grpc

import (
	"context"
	"fmt"

	"github.com/nolafw/projecttemplate/internal/module/user/service"
	"github.com/nolafw/projecttemplate/internal/plamo/dikit"
	pb "github.com/nolafw/projecttemplate/service/adapter/user"
	"google.golang.org/grpc"
)

// gRPCでの接続処理
type UserGRPCService struct {
	pb.UnimplementedUserServer
	dikit.GRPCServiceRegistrar
	service service.UserService
}

// DIで使う必要なserviceをinject
func NewUserGRPCService(service service.UserService) *UserGRPCService {
	return &UserGRPCService{
		service: service,
	}
}

// IMPORTANT! gRPCサーバーに登録するためのメソッド
func (s *UserGRPCService) RegisterWithServer(grpcServer *grpc.Server) {
	pb.RegisterUserServer(grpcServer, s)
}

// TODO:
// serviceは基本的に、こんな感じの実装をするルールにする
// context.Contextを引数に受け取る
// reqのところは、DTOで置き換える
// 返り値は、DTOとerror
func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	// panic("just for test") // リカバリのテスト用
	fmt.Println("GetUser called with request:", req.UserId)

	return &pb.GetUserResponse{
		UserId: "1",
		Name:   "John Doe",
		Email:  "j@example.com",
	}, nil
}
