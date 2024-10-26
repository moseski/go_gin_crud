package controllers

import (
    "context"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/moseski/go_gin_crud/initializers"
    "github.com/moseski/go_gin_crud/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "net/http"
    "time"
)

var ctx = context.Background()

// CreateUser - Create a new user in MongoDB
func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = user.CreatedAt

    collection := initializers.MongoDB.Collection("users")
    _, err := collection.InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Clear cache for all users
    initializers.RedisClient.Del(ctx, "all_users")

    c.JSON(http.StatusCreated, user)
}

// GetUser - Get a single user by ID from MongoDB with Redis caching
func GetUser(c *gin.Context) {
    userID := c.Param("id")
    cacheKey := "user:" + userID

    // Try to get cached data
    cachedUser, err := initializers.RedisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        var user models.User
        if json.Unmarshal([]byte(cachedUser), &user) == nil {
            c.JSON(http.StatusOK, user)
            return
        }
    }

    // Cache miss, get from database
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    var user models.User
    collection := initializers.MongoDB.Collection("users")
    err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Cache the response
    userJSON, _ := json.Marshal(user)
    initializers.RedisClient.Set(ctx, cacheKey, userJSON, 5*time.Minute)

    c.JSON(http.StatusOK, user)
}

// GetAllUsers - Get all users with Redis caching
func GetAllUsers(c *gin.Context) {
    cacheKey := "all_users"

    cachedUsers, err := initializers.RedisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        var users []models.User
        if json.Unmarshal([]byte(cachedUsers), &users) == nil {
            c.JSON(http.StatusOK, users)
            return
        }
    }

    var users []models.User
    collection := initializers.MongoDB.Collection("users")
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user"})
            return
        }
        users = append(users, user)
    }

    if cursor.Err() != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating cursor"})
        return
    }

    // Cache the result
    usersJSON, _ := json.Marshal(users)
    initializers.RedisClient.Set(ctx, cacheKey, usersJSON, 5*time.Minute)

    c.JSON(http.StatusOK, users)
}

// UpdateUser - Update user and clear relevant cache entries
func UpdateUser(c *gin.Context) {
    userID := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    var updatedUser models.User
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    update := bson.M{
        "$set": bson.M{
            "name":      updatedUser.Name,
            "email":     updatedUser.Email,
            "password":  updatedUser.Password,
            "updatedAt": time.Now(),
        },
    }

    collection := initializers.MongoDB.Collection("users")
    _, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    // Clear cache for the updated user and all users list
    initializers.RedisClient.Del(ctx, "user:"+userID, "all_users")

    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser - Delete user and clear relevant cache entries
func DeleteUser(c *gin.Context) {
    userID := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    collection := initializers.MongoDB.Collection("users")
    _, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    // Clear cache for the deleted user and all users list
    initializers.RedisClient.Del(ctx, "user:"+userID, "all_users")

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
