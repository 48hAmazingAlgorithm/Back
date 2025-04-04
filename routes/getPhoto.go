package routes

import (
	
	"context"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetClient(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Client")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var clients []Client
	_ = cursor.All(ctx, &clients)
	c.JSON(http.StatusOK, gin.H{
		"message": clients,
	})
}