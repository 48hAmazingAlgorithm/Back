package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	//"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Client struct {
	NumClient string `bson:"_id,omitempty" json:"numero_client"`
	NomClient string `bson:"nom_client" json:"nom_client"`
}


var Mongoclient *mongo.Client

func connectMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb+srv://challenge48hynov:cha!!enge48ynov@challenge48h.xgbo0vx.mongodb.net/?retryWrites=true&w=majority&appName=Challenge48h")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Impossible de se connecter à MongoDB")
	}

	fmt.Println("Connexion MongoDB réussie")
}

func GetClient(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Client")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel();
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var clients []Client
	_ = cursor.All(ctx, &clients)
	c.JSON(http.StatusOK, gin.H{
		"message" : clients,
	})
}

func PostClient(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Client")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel();
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

func main() {
	connectMongoDB()

	router := gin.Default();

	router.GET("/get", GetClient)
	router.POST("/post", PostClient)
	router.Run(":8080")
}

