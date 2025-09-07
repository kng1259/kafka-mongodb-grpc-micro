package main

import (
	"context"
	"fmt"
	"kmgm-consumer/config"
	"kmgm-consumer/grpcServer"
	"kmgm-consumer/kafka"
	"kmgm-consumer/models"
	"kmgm-consumer/mongodb"
	product "kmgm-consumer/protogen"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Check MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	// Get product collection and create repository
	collection := client.Database(cfg.MongoDB.Database).Collection(cfg.MongoDB.Collection)
	repo := mongodb.NewRepository(collection)

	// Initialize Kafka consumer
	if cfg.Kafka.Enabled {
		// Start gRPC server in a goroutine
		go startGRPCServer(repo, cfg.GRPC.Port)
		kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
		defer kafkaConsumer.Close()

		log.Println("Starting consumer...")

		err = kafkaConsumer.Consume(context.Background(), func(p models.Product) error {
			return repo.StoreProduct(context.Background(), p)
		})
		if err != nil {
			log.Fatalf("Error in consumer: %v", err)
		}
	} else {
		// Start gRPC server
		startGRPCServer(repo, cfg.GRPC.Port)
	}
}

func startGRPCServer(repo *mongodb.Repository, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	grpcServer := grpcServer.NewServer(repo)

	// Register the service
	product.RegisterProductServiceServer(s, grpcServer)

	log.Printf("gRPC server listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
