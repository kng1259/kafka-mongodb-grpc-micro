package grpcServer

import (
	"context"
	"kmgm-consumer/models"
	"kmgm-consumer/mongodb" // Add this import
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	product "kmgm-consumer/protogen"
)

type Server struct {
	product.UnimplementedProductServiceServer
	repo *mongodb.Repository // Change from collection to repo
}

func NewServer(repo *mongodb.Repository) *Server { // Update constructor
	return &Server{repo: repo}
}

func (s *Server) GetProducts(ctx context.Context, req *product.GetProductsRequest) (*product.GetProductsResponse, error) {
	// Build filter
	filter := bson.M{}
	if req.Category != "" {
		filter["category"] = req.Category
	}

	// Set pagination
	page := int64(req.Page)
	if page <= 0 {
		page = 1
	}
	limit := int64(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	skip := (page - 1) * limit

	// Get total count
	total, err := s.repo.GetCollection().CountDocuments(ctx, filter) // Use repo's collection
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count products: %v", err)
	}

	// Find products
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.repo.GetCollection().Find(ctx, filter, findOptions) // Use repo's collection
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find products: %v", err)
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	for cursor.Next(ctx) {
		var p models.Product
		if err := cursor.Decode(&p); err != nil {
			log.Printf("Failed to decode product: %v", err)
			continue
		}

		products = append(products, &product.Product{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			Category:    p.Category,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "cursor error: %v", err)
	}

	return &product.GetProductsResponse{
		Products: products,
		Total:    int32(total),
		Page:     int32(page),
		Limit:    int32(limit),
	}, nil
}

func (s *Server) GetProductByID(ctx context.Context, req *product.GetProductByIDRequest) (*product.GetProductByIDResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	filter := bson.M{"_id": objectID}
	var p models.Product

	err = s.repo.GetCollection().FindOne(ctx, filter).Decode(&p) // Use repo's collection
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find product: %v", err)
	}

	return &product.GetProductByIDResponse{
		Product: &product.Product{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			Category:    p.Category,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}
