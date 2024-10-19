package main

import (
    "github.com/gin-gonic/gin"
    "github.com/moseski/go_gin_crud/initializers"
    "github.com/moseski/go_gin_crud/routes"
)

func init() {
    // Load environment variables
    initializers.LoadEnvVariables()

    // Connect to the database
    initializers.ConnectToDB()
}

func main() {
    // Use the existing routing setup
    r := gin.Default()

    // Setup user routes
    routes.UserRoutes(r)

    // Run server on port 8080
    r.Run(":8080")
}


