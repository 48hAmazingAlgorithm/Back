package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	PhotoID            string    `bson:"photo_id"`
	DatePhoto          time.Time `bson:"date_photo"`
}

var Mongoclient *mongo.Client
var bucket *gridfs.Bucket
var err = godotenv.Load()
var encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))

func connectMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb+srv://challenge48hynov:cha!!enge48ynov@challenge48h.xgbo0vx.mongodb.net/?retryWrites=true&w=majority&appName=Challenge48h")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Mongoclient, err = mongo.Connect(context.Background(),clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = Mongoclient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Impossible de se connecter à MongoDB")
	}
	bucket, _ = gridfs.NewBucket(Mongoclient.Database("Challenge48h"))
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
	err := collectionC.FindOne(ctx, bson.M{"_id": individu.NumeroClient}).Decode(&existingClient)
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

func PostPhoto(c *gin.Context) {
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucune photo envoyé"})
		return
	}
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	var img image.Image
	if fileExt == ".png" {
		img, err = png.Decode(file)
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		img, err = jpeg.Decode(file)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format non supporté (JPG, PNG uniquement)"})
		return
	}
	img = addFilligrane(img)
	var buf bytes.Buffer
	if fileExt == ".png" {
		err = png.Encode(&buf, img)
	} else {
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur d'encodage de l'image"})
		return
	}


	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var idIndividu struct {
		Id string `json:"id_individu"`
	}
	if err := c.ShouldBindJSON(&idIndividu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
		})
		return
	}
	var individu Individu
	err = collection.FindOne(ctx, bson.M{"_id" : idIndividu.Id}).Decode(&individu)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "L'individu n'existe pas",
		})
		return
	}

	newId, err :=uploadFile(buf)
	newId, _ = encryptID(newId)
	individu.PhotoID = newId
	individu.DatePhoto = time.Now()
	result, err := collection.UpdateOne(ctx,bson.M{"_id" : idIndividu.Id},individu)
	if err != nil {
		log.Fatal("probleme lors de l'update",err)
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func addFilligrane(img image.Image) image.Image{
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
	return context.Image()
}

func uploadFile(buf bytes.Buffer) (string, error) {
	name := uuid.New().String()
	uploadStream, _ := bucket.OpenUploadStream(name)
	defer uploadStream.Close()
	uploadStream.Write(buf.Bytes())
	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

func encryptID(text string) (string, error) {
	block, _ := aes.NewCipher(encryptionKey)
	nonce := make([]byte,12)
	io.ReadFull(rand.Reader, nonce)
	aesGCM, _ :=cipher.NewGCM(block)
	ciphertext := aesGCM.Seal(nil, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func decryptID(text string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(text)
	nonce := data[:12]         	
	ciphertext := data[12:]    
	block, _ := aes.NewCipher(encryptionKey)
	aesGCM, _ := cipher.NewGCM(block)
	plainText, _ := aesGCM.Open(nil, nonce, ciphertext, nil)
	return string(plainText), nil
}


func main() {
	connectMongoDB()
	router := gin.Default()

	router.GET("/getClient", GetClient)
	router.POST("/postClient", PostClient)
	router.GET("/getIndividus", GetIndividus)
	router.POST("/postIndividu", PostIndividu)
	router.Run(":8080")
}
