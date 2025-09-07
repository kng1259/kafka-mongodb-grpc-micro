package grpcClient

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	product "kmgm-producer/protogen"
)

type Client struct {
	conn    *grpc.ClientConn
	service product.ProductServiceClient
}

func NewClient(host string, port int, timeout time.Duration) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %v", addr, err)
	}

	log.Printf("Connected to gRPC server at %s", addr)

	return &Client{
		conn:    conn,
		service: product.NewProductServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetProducts(page, limit int32, category string) (*product.GetProductsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.service.GetProducts(ctx, &product.GetProductsRequest{
		Page:     page,
		Limit:    limit,
		Category: category,
	})
}

func (c *Client) GetProductByID(id string) (*product.GetProductByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.service.GetProductByID(ctx, &product.GetProductByIDRequest{
		Id: id,
	})
}
