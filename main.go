package main

import (
    "github.com/gin-gonic/gin"
    "github.com/moseski/go_gin_crud/initializers"
    "github.com/moseski/go_gin_crud/routes"
)

func init() {
    // Load environment variables and initialize MongoDB and Redis connections
    initializers.LoadEnvVariables()
    initializers.ConnectToMongoDB()
    initializers.InitRedis()  // Initialize Redis connection
}

func main() {
    r := gin.Default()

    // Register routes
    routes.UserRoutes(r)

    // Start the server on port 8080
    r.Run(":8080")
}