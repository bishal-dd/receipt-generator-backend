package main

import (
	"log"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	resolver "github.com/bishal-dd/receipt-generator-backend/graph/resolver"
	"github.com/bishal-dd/receipt-generator-backend/pkg/db"
	"github.com/bishal-dd/receipt-generator-backend/pkg/redis"
	"github.com/bishal-dd/receipt-generator-backend/pkg/rmq"
	"github.com/bishal-dd/receipt-generator-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	 err := godotenv.Load()
	 if err != nil {
	   log.Fatal("Error loading .env file")
	 }
	 database := db.Init()
	 cacheRedis, queueRedis, err := redis.Init()
	 if err != nil {
		log.Fatal(err)
	 }
	 if err := rmq.InitEmailQueue(queueRedis); err != nil {
		log.Fatal(err)
	 }
	 dependencyResolver := resolver.InitializeResolver(cacheRedis, database)
 
	// Setting up Gin
	log.Printf("connect to http://localhost:%d/graphql for GraphQL playground", 8080)
	r := gin.Default()
	r.GET("/graphql", routes.PlaygroundHandler())
	r.Use(GinContextToContextMiddleware())
	r.Use(loaders.LoaderMiddleware(database))
	r.POST("/query", routes.GraphqlHandler(dependencyResolver))
	r.GET("/issuePresignedURL", routes.HandlePresignedURL)
	r.Run()
	
}