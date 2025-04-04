package main

import (
	"bytes"
	"context"

	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	routes "challenge/routes"
)

type Client struct {
	NumClient string `bson:"_id,omitempty" json:"numero_client"`
	NomClient string `bson:"nom_client" json:"nom_client"`
}

type Individu struct {
	IDIndividu         primitive.ObjectID    `bson:"_id,omitempty" json:"id_individu"`
	Nom                string    `bson:"nom" json:"nom"`
	Prenom             string    `bson:"prenom" json:"prenom"`
	DateNaissance      time.Time `bson:"date_naissance" json:"date_naissance"`
	DateFinValiditeCNI time.Time `bson:"date_fin_validite_CNI" json:"date_fin_validite_CNI"`
	NumeroCNI          string    `bson:"numero_CNI" json:"numero_CNI"`
	PhotoID            string    `bson:"photo_id"`
	DatePhoto          time.Time `bson:"date_photo"`
}

var Mongoclient *mongo.Client
var bucket *gridfs.Bucket


func connectMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb+srv://challenge48hynov:cha!!enge48ynov@challenge48h.xgbo0vx.mongodb.net/?retryWrites=true&w=majority&appName=Challenge48h")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Mongoclient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = Mongoclient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Impossible de se connecter à MongoDB")
	}
	bucket, _ = gridfs.NewBucket(Mongoclient.Database("Challenge48h"), options.GridFSBucket().SetName("images"))
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

func PostPhoto(c *gin.Context) {
	var requestData struct {
		Id_individu string `json:"id_individu"`
		Photo_data string `json:"photo"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
		})
		return
	}
	photoData, err := base64.StdEncoding.DecodeString(requestData.Photo_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de décodage de l'image"})
		return
	}
	
	img, tipe, err := image.Decode(bytes.NewReader(photoData))
	log.Println("type de l'image", tipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de décodage de l'image"})
		return
	}
	img = addFilligrane(img)
	newId, err := uploadFile(img)
	newId, _ = routes.EncryptID(newId)

	idObjectId, err := primitive.ObjectIDFromHex(requestData.Id_individu)
	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var individu Individu
	err = collection.FindOne(ctx, bson.M{"_id": idObjectId}).Decode(&individu)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "L'individu n'existe pas",
		})
		return
	}

	individu.PhotoID = newId
	individu.DatePhoto = time.Now()
	result, err := collection.UpdateOne(ctx, bson.M{"_id": idObjectId}, bson.M{"$set": bson.M{"photo_id": individu.PhotoID,"date_photo": individu.DatePhoto,}})
	if err != nil {
		log.Fatal("probleme lors de l'update", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func GetPhoto(c *gin.Context) {
	clientId := c.Param("client_id")
	if clientId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client ID is required"})
		return
	}
	log.Println("clientId reçu :", clientId)  

	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var individu Individu
	idObjectId, err := primitive.ObjectIDFromHex(clientId)
	if err != nil {
		log.Println("Erreur de conversion de l'ID du client :", clientId)  
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID du client invalide"})
		return
	}

	err = collection.FindOne(ctx, bson.M{"_id": idObjectId}).Decode(&individu)
	if err != nil {
		log.Println("Erreur lors de la recherche de l'individu :", err)  
		c.JSON(http.StatusNotFound, gin.H{"error": "Client non trouvé"})
		return
	}

	
	log.Println("Individu trouvé :", individu)

	
	photoId, err := routes.DecryptID(individu.PhotoID)
	if err != nil {
		log.Println("Erreur lors du décryptage de l'ID de la photo :", err)  
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du décryptage de l'ID de la photo"})
		return
	}

	log.Println("photoId décrypté :", photoId)
	
	objectId, err := primitive.ObjectIDFromHex(photoId)
	if err != nil {
		log.Println("Erreur lors de la conversion de l'ID de la photo en ObjectId :", err)  
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de la photo invalide"})
		return
	}
	log.Println("id de la photo", objectId)
	
	downloadStream, err := bucket.OpenDownloadStream(objectId)
	if err != nil {
		log.Println("Erreur lors du téléchargement de la photo depuis GridFS :", err)  
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du téléchargement de la photo"})
		return
	}
	defer downloadStream.Close()

	var photoData bytes.Buffer
	_, err = photoData.ReadFrom(downloadStream)
	if err != nil {
		
		log.Println("Erreur lors de la lecture des données de l'image :", err)  
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des données de l'image"})
		return
	}
	photo, _ , err := image.Decode(&photoData)
	if err != nil {
		log.Println("errer lors de image.Decode", err)
		return
	}
	outfile, err := os.Create("getPhoto.png")
	if err != nil {
		log.Println("error lors de la creation de fichier", err)
		return
	}

	png.Encode(outfile,photo)
	outfile.Close()

	var imgBuffer bytes.Buffer
	c.Header("Content-Type", "image/png")
	png.Encode(&imgBuffer, photo)
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), imgBuffer.Bytes())
		
	log.Println("Image envoyée avec succès.")
}

func addFilligrane(img image.Image) image.Image {
    if img == nil {
        log.Fatal("Erreur : l'image est nil")
        return nil
    }

    context := gg.NewContextForImage(img)
    if context == nil {
        log.Fatal("Erreur : impossible de créer le contexte pour l'image")
        return nil
    }

    imgWidth := img.Bounds().Max.X
    imgHeight := img.Bounds().Max.Y
    context.SetRGBA(0, 0, 0, 1)
    fontSize := float64(imgHeight) * 0.1

    err := context.LoadFontFace("c:/Windows/Fonts/Amiri-Bold.ttf", fontSize)
    if err != nil {
        log.Fatal("Erreur lors du chargement de la police: ", err)
        return nil
    }

    text := "CHALLENGE 48H YNOV"
    centerX := float64(imgWidth) / 2
    centerY := float64(imgHeight) / 2
    angle := 45.0

    context.RotateAbout(gg.Radians(angle), centerX, centerY)
    context.DrawStringAnchored(text, centerX, centerY, 0.5, 0.5)

	context.SavePNG("image.png")
    return context.Image()
}

func uploadFile(img image.Image) (string, error) {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	name := uuid.New().String()
	uploadStream, _ := bucket.OpenUploadStream(name)
	defer uploadStream.Close()
	uploadStream.Write(buf.Bytes())
	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

func getAllDatabase() {
	collection := Mongoclient.Database("Challenge48h").Collection("fs.files")
	log.Println("database : ", collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var files []struct {
		Id string `bson:"_id"`
		Length int `bson:"length"`
		ChunkSize int `bson:"chunkSize"`
		UploadDate time.Time `bson:"uploadDate"`
		FileName string `bson:"filename"`
	}
	_ = cursor.All(ctx, &files)
	log.Println(files)
	for _, file := range files {
		log.Println("file : ",file)
	}
	
}

func main() {
	connectMongoDB()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Autoriser le frontend sur localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	getAllDatabase()
	router.GET("/getClient", GetClient)
	router.POST("/postClient", PostClient)
	router.GET("/getIndividus", GetIndividus)
	router.POST("/postIndividu", PostIndividu)
	router.GET("/getPhoto/:client_id", GetPhoto)
	router.POST("/postPhoto", PostPhoto)

	router.Run(":8080")
}
