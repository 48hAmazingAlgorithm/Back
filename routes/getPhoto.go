package routes

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	
	photoId, err := DecryptID(individu.PhotoID)
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
	
	downloadStream, err := Bucket.OpenDownloadStream(objectId)
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

