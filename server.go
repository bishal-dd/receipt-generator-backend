package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bishal-dd/receipt-generator-backend/graph"
	resolver "github.com/bishal-dd/receipt-generator-backend/graph/resolver"
	"github.com/bishal-dd/receipt-generator-backend/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	 db.Init()
	 err := godotenv.Load()
	 if err != nil {
	   log.Fatal("Error loading .env file")
	 }
	// Setting up Gin
	log.Printf("connect to http://localhost:%d/graphql for GraphQL playground", 8080)
	r := gin.Default()
	r.Use(GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler())
	r.GET("/graphql", playgroundHandler())
	r.Run()
}