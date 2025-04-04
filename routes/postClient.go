package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PostClient(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Client")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var client Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
		})
		return
	}
	result, err := collection.InsertOne(ctx, client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de l'insertion en base",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}