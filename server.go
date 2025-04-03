package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"net/http"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Client struct {
	NumClient string `bson:"_id,omitempty" json:"numero_client"`
	NomClient string `bson:"nom_client" json:"nom_client"`
}

type Individu struct {
	IDIndividu         string    `bson:"_id,omitempty" json:"id_individu"`
	Nom                string    `bson:"nom" json:"nom"`
	Prenom             string    `bson:"prenom" json:"prenom"`
	DateNaissance      time.Time `bson:"date_naissance" json:"date_naissance"`
	DateFinValiditeCNI time.Time `bson:"date_fin_validite_CNI" json:"date_fin_validite_CNI"`
	NumeroCNI          string    `bson:"numero_CNI" json:"numero_CNI"`
	NumeroClient       string    `bson:"numero_client" json:"numero_client"`
}

var Mongoclient *mongo.Client

func connectMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb+srv://challenge48hynov:cha!!enge48ynov@challenge48h.xgbo0vx.mongodb.net/?retryWrites=true&w=majority&appName=Challenge48h")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Mongoclient, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = Mongoclient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Impossible de se connecter à MongoDB")
	}

	fmt.Println("Connexion MongoDB réussie")
}

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

func PostIndividu(c *gin.Context) {
	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var individu Individu
	if err := c.ShouldBindJSON(&individu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
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
		"message": result,
	})
}

func addFilligrane(img image.Image) {
	context := gg.NewContextForImage(img)

	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y
	context.SetRGBA(0, 0, 0, 1)
	fontSize := float64(imgHeight) * 0.1
	err := context.LoadFontFace("c:/Windows/Fonts/Amiri-Bold.ttf", fontSize)
	if err != nil {
		log.Fatal("Erreur lors du chargement de la police: ", err)
	}
	text := "CHALLENGE 48H YNOV"
	centerX := float64(imgWidth) / 2
	centerY := float64(imgHeight) / 2
	angle := 45.0

	context.RotateAbout(gg.Radians(angle), centerX, centerY)

	context.DrawStringAnchored(text, centerX, centerY, 0.5, 0.5)

	err = context.SavePNG("./imageFinale.png")
	if err != nil {
		log.Fatal("Erreur lors de la sauvegarde de l'image: ", err)
	}
}

func main() {
	imgPath := "./image.png"
	img, err := gg.LoadImage(imgPath)
	if err != nil {
		log.Fatal("Erreur lors du chargement de l'image: ", err)
	}
	addFilligrane(img)
	connectMongoDB()
	router := gin.Default()

	router.GET("/getClient", GetClient)
	router.POST("/postClient", PostClient)
	router.GET("/getIndividus", GetIndividus)
	router.POST("/postIndividu", PostIndividu)
	router.Run(":8080")
}
