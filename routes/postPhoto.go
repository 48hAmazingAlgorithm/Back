package routes

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostPhotoRecto(c *gin.Context) {
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
	img, IDIndividu := getImage(requestData, c)
	img = addFilligrane(img)
	newId, err := uploadFile(img)
	newId, _ = EncryptID(newId)

	individu, collection, ctx, ObjectIdIndividu, cancel := GetSingleIndividu(IDIndividu, c)
	defer cancel()
	individu.PhotoRectoID = newId
	individu.DatePhoto = time.Now()
	result, err := collection.UpdateOne(ctx, bson.M{"_id": ObjectIdIndividu}, bson.M{"$set": bson.M{"photoRecto_id": individu.PhotoRectoID,"date_photo": individu.DatePhoto,}})
	if err != nil {
		log.Fatal("probleme lors de l'update", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func PostPhotoVerso(c *gin.Context) {
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
	img, IDIndividu := getImage(requestData, c)
	img = addFilligrane(img)
	newId, err := uploadFile(img)
	newId, _ = EncryptID(newId)

	individu, collection, ctx, ObjectIdIndividu, cancel := GetSingleIndividu(IDIndividu, c)
	defer cancel()

	individu.PhotoVersoID = newId
	individu.DatePhoto = time.Now()
	result, err := collection.UpdateOne(ctx, bson.M{"_id": ObjectIdIndividu}, bson.M{"$set": bson.M{"photoVerso_id": individu.PhotoVersoID,"date_photo": individu.DatePhoto,}})
	if err != nil {
		log.Fatal("probleme lors de l'update", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
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

    context.SetRGBA(0, 0, 0, 1)
    fontSize := float64(img.Bounds().Max.Y) * 0.1
    err := context.LoadFontFace("c:/Windows/Fonts/Amiri-Bold.ttf", fontSize)
    if err != nil {
        log.Fatal("Erreur lors du chargement de la police: ", err)
        return nil
    }

    text := "CHALLENGE 48H YNOV"
    centerX := float64(img.Bounds().Max.X) / 2
    centerY := float64(img.Bounds().Max.Y) / 2
    angle := 45.0

    context.RotateAbout(gg.Radians(angle), centerX, centerY)
    context.DrawStringAnchored(text, centerX, centerY, 0.5, 0.5)
    return context.Image()
}

func uploadFile(img image.Image) (string, error) {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	name := uuid.New().String()
	uploadStream, _ := Bucket.OpenUploadStream(name)
	defer uploadStream.Close()
	uploadStream.Write(buf.Bytes())
	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

func getImage(requestData struct {
	Id_individu string `json:"id_individu"`;Photo_data string `json:"photo"`}, c *gin.Context) (image.Image, string) {

	photoData, err := base64.StdEncoding.DecodeString(requestData.Photo_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de décodage de l'image"})
		return nil, ""
	}
	
	img, tipe, err := image.Decode(bytes.NewReader(photoData))
	log.Println("type de l'image", tipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de décodage de l'image"})
		return nil, ""
	}
	return img, requestData.Id_individu
}

func GetSingleIndividu(id string, c *gin.Context) (Individu, *mongo.Collection, context.Context, primitive.ObjectID, context.CancelFunc) {
	idObjectId, err := primitive.ObjectIDFromHex(id)
	collection := Mongoclient.Database("Challenge48h").Collection("Individu")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	var individu Individu
	err = collection.FindOne(ctx, bson.M{"_id": idObjectId}).Decode(&individu)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "L'individu n'existe pas",
		})
		return individu, nil, nil, primitive.NilObjectID, nil
	}
	return individu, collection, ctx, idObjectId, cancel
}