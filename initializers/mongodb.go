package initializers

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

// ConnectToMongoDB connects to MongoDB and sets the global MongoDB variable
func ConnectToMongoDB() {
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("MONGODB_URI not found in environment variables")
    }
    clientOptions := options.Client().ApplyURI(mongoURI)

    // Connect to MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping the MongoDB server
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    // Set the MongoDB variable to use the "users" database
    MongoDB = client.Database("users")

    log.Println("Connected to MongoDB!")
}