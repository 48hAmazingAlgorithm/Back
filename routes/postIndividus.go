package routes

import (
	"net/http"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)



func PostIndividu(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	collectionC := Mongoclient.Database("Challenge48h").Collection("Client")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var individu Individu
	if err := c.ShouldBindJSON(&individu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
		})
		return
	}

	var existingClient Client
	err := collectionC.FindOne(ctx, bson.M{"_id": individu.IDIndividu}).Decode(&existingClient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Le client n'existe pas",
		})
		return
	}

	result, err := collection.InsertOne(ctx, individu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de l'insertion en base",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result.InsertedID,
	})
}