package routes

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetIndividu(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(ctx)
	var individus []Individu
	err = cursor.All(ctx, &individus)
	if err != nil {
		log.Println(err)
	}
	log.Println(individus)
	c.JSON(http.StatusOK, individus)
}
