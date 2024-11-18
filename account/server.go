package account

import (
	"context"
	"fmt"
	"net"

	pb "github.com/silven-dynamics/go-ecommerce/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	accountService AccountService
}

func ListenGRPC(s AccountService, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		accountService: s,
	})
	reflection.Register(serv)
	return serv.Serve(listener)
}

func (s *grpcServer) PostAccount(
	ctx context.Context,
	r *pb.PostAccountRequest,
) (*pb.PostAccountResponse, error) {
	account, err := s.accountService.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   account.ID,
			Name: account.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(
	ctx context.Context,
	r *pb.GetAccountRequest,
) (*pb.GetAccountResponse, error) {
	account, err := s.accountService.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   account.ID,
			Name: account.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(
	ctx context.Context,
	r *pb.GetAccountsRequest,
) (*pb.GetAccountsResponse, error) {
	res, err := s.accountService.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, p := range res {
		accounts = append(accounts,
			&pb.Account{
				Id:   p.ID,
				Name: p.Name,
			},
		)
	}

	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
