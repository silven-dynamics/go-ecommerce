package catalog

import (
	"context"
	"fmt"
	"net"

	pb "github.com/stiffinWanjohi/go-ecommerce/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	catalogService CatalogService
}

func ListenGRPC(s CatalogService, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{
		catalogService: s,
	})
	reflection.Register(serv)
	return serv.Serve(listener)
}

func (s *grpcServer) PostProduct(
	ctx context.Context,
	r *pb.PostProductRequest,
) (*pb.PostProductResponse, error) {
	product := Product{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
	}
	p, err := s.catalogService.PostProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProduct(
	ctx context.Context,
	r *pb.GetProductRequest,
) (*pb.GetProductResponse, error) {
	p, err := s.catalogService.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(
	ctx context.Context,
	r *pb.GetProductsRequest,
) (*pb.GetProductsResponse, error) {
	var products []Product
	var err error
	if len(r.Ids) > 0 {
		products, err = s.catalogService.GetProductsByIDs(ctx, r.Ids)
	} else if r.Query != "" {
		products, err = s.catalogService.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else {
		products, err = s.catalogService.GetProducts(ctx, r.Skip, r.Take)
	}

	if err != nil {
		return nil, err
	}

	pbProducts := make([]*pb.Product, len(products))
	for index, p := range products {
		pbProducts[index] = &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}

	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}
