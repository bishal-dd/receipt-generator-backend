package main

import (
	"log"
	"os"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	resolver "github.com/bishal-dd/receipt-generator-backend/graph/resolver"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/pkg/crypto"
	"github.com/bishal-dd/receipt-generator-backend/pkg/db"
	"github.com/bishal-dd/receipt-generator-backend/pkg/redis"
	"github.com/bishal-dd/receipt-generator-backend/pkg/rmq"
	"github.com/bishal-dd/receipt-generator-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	env := os.Getenv("ENV")
	if env == "development" {
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	privateKey, err := crypto.LoadPrivateKey()
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}

	publicKey, err := crypto.LoadPublicKey()
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}
	encryption.InitKeys(publicKey, privateKey)
	log.Println("RSA keys loaded successfully")

	publicKeyPEM, err := crypto.PublicKeyToPEM(publicKey)
	if err != nil {
		log.Fatalf("failed to convert public key to PEM: %v", err)
	}
	privateKeyPEM, err := crypto.PrivateKeyToPEM(privateKey)
	if err != nil {
		log.Fatalf("failed to convert private key to PEM: %v", err)
	}

	InitializeApi()
	database := db.Init()
	queueRedis, err := redis.Init()
	if err != nil {
		log.Fatal(err)
	}
	if err := rmq.InitEmailQueue(queueRedis); err != nil {
		log.Fatal(err)
	}

	httpClient := resty.New()

	dependencyResolver := resolver.InitializeResolver(database, httpClient, publicKeyPEM, privateKeyPEM)

	log.Printf("connect to http://localhost:%d/graphql for GraphQL playground", 8080)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "https://2quickpaper.com", "https://www.2quickpaper.com", "https://www.billstotrack.com", "https://billstotrack.com"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/graphql", routes.PlaygroundHandler())
	r.POST("/profile", func(c *gin.Context) {
		routes.AddProfile(c, database)
	})
	r.POST("/payment-webhook", func(c *gin.Context) {
		routes.PaymentWebhook(c, database)
	})
	r.GET("/issuePresignedURL", routes.HandlePresignedURL)
	r.Use(AuthMiddleware())
	r.Use(loaders.LoaderMiddleware(database))
	r.POST("/query", routes.GraphqlHandler(dependencyResolver))
	r.Run()

}
