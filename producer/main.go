package main

import (
	"context"
	"encoding/json"
	"kmgm-producer/config"
	"kmgm-producer/grpcClient"
	"kmgm-producer/kafka"
	"kmgm-producer/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Kafka producer
	var kafkaProducer *kafka.Producer
	var err error
	if cfg.Kafka.Enabled {
		kafkaProducer, err = kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Username, cfg.Kafka.Password, cfg.Kafka.Topic)
		if err != nil {
			log.Fatalf("Failed to create Kafka producer: %v", err)
		}
		defer kafkaProducer.Close()
	}

	// Set up Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Product creation endpoint
	if cfg.Kafka.Enabled {
		r.POST("/"+cfg.Server.Route, func(c *gin.Context) {
			var product models.Product
			if err := c.ShouldBindJSON(&product); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Set creation timestamp
			product.CreatedAt = time.Now()

			// Convert product to JSON
			productJSON, err := json.Marshal(product)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal product"})
				return
			}

			// Produce message to Kafka
			err = kafkaProducer.ProduceMessage(context.Background(), []byte(product.Name), productJSON)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write message to Kafka: " + err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Product sent to Kafka",
				"product": product,
			})
		})
	}

	// Initialize gRPC client
	grpcClient, err := grpcClient.NewClient(
		cfg.GRPC.Host,
		cfg.GRPC.Port,
		cfg.GRPC.Timeout,
	)
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcClient.Close()

	// Get products endpoint (uses gRPC to consumer)
	r.GET("/"+cfg.Server.Route, func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		category := c.Query("category")

		response, err := grpcClient.GetProducts(int32(page), int32(limit), category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": response.Products,
			"total":    response.Total,
			"page":     response.Page,
			"limit":    response.Limit,
		})
	})

	// Get product by ID endpoint
	r.GET("/"+cfg.Server.Route+"/:id", func(c *gin.Context) {
		id := c.Param("id")

		response, err := grpcClient.GetProductByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"product": response.Product,
		})
	})

	log.Printf("Producer service starting on port %d", cfg.Server.Port)
	log.Fatal(r.Run(":" + viper.GetString("server.port")))
}
