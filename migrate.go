package main

import (
    "github.com/moseski/go_gin_crud/initializers"
    "github.com/moseski/go_gin_crud/models"
)

func init() {
    initializers.LoadEnvVariables()  // Load environment variables (e.g., database URL)
    initializers.ConnectToDB()       // Connect to the database
}

func main() {
    initializers.DB.AutoMigrate(&models.User{})  // Migrate the User model
}