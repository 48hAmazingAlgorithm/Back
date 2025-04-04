package routes

import(
	"time"
	"context"
	"github.com/gin-gonic/gin"

	"net/http"
	"go.mongodb.org/mongo-driver/bson"
)

func GetIndividus(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var individus []Individu
	_ = cursor.All(ctx, &individus)
	c.JSON(http.StatusOK, gin.H{
		"message": individus,
	})
}