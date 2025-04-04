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

func GetPhotoRecto(c *gin.Context) {
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

	photoRectoId, err := DecryptID(individu.PhotoRectoID)
	if err != nil {
		log.Println("Erreur lors du décryptage de l'ID de la photoRecto :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du décryptage de l'ID de la photoRecto"})
		return
	}

	log.Println("photoRectoId décrypté :", photoRectoId)

	objectId, err := primitive.ObjectIDFromHex(photoRectoId)
	if err != nil {
		log.Println("Erreur lors de la conversion de l'ID de la photoRecto en ObjectId :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de la photoRecto invalide"})
		return
	}
	log.Println("id de la photoRecto", objectId)

	downloadStream, err := Bucket.OpenDownloadStream(objectId)
	if err != nil {
		log.Println("Erreur lors du téléchargement de la photoRecto depuis GridFS :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du téléchargement de la photoRecto"})
		return
	}
	defer downloadStream.Close()

	var photoRectoData bytes.Buffer
	_, err = photoRectoData.ReadFrom(downloadStream)
	if err != nil {

		log.Println("Erreur lors de la lecture des données de l'image :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des données de l'image"})
		return
	}
	photoRecto, _, err := image.Decode(&photoRectoData)
	if err != nil {
		log.Println("errer lors de image.Decode", err)
		return
	}
	outfile, err := os.Create("getPhotoRecto.png")
	if err != nil {
		log.Println("error lors de la creation de fichier", err)
		return
	}

	png.Encode(outfile, photoRecto)
	outfile.Close()

	var imgBuffer bytes.Buffer
	c.Header("Content-Type", "image/png")
	png.Encode(&imgBuffer, photoRecto)
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), imgBuffer.Bytes())

	log.Println("Image envoyée avec succès.")
}


func GetPhotoVerso(c *gin.Context) {
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

	photoVersoId, err := DecryptID(individu.PhotoVersoID)
	if err != nil {
		log.Println("Erreur lors du décryptage de l'ID de la photoVerso :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du décryptage de l'ID de la photoVerso"})
		return
	}

	log.Println("photoVersoId décrypté :", photoVersoId)

	objectId, err := primitive.ObjectIDFromHex(photoVersoId)
	if err != nil {
		log.Println("Erreur lors de la conversion de l'ID de la photoVerso en ObjectId :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de la photoVerso invalide"})
		return
	}
	log.Println("id de la photoVerso", objectId)

	downloadStream, err := Bucket.OpenDownloadStream(objectId)
	if err != nil {
		log.Println("Erreur lors du téléchargement de la photoVerso depuis GridFS :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du téléchargement de la photoVerso"})
		return
	}
	defer downloadStream.Close()

	var photoVersoData bytes.Buffer
	_, err = photoVersoData.ReadFrom(downloadStream)
	if err != nil {

		log.Println("Erreur lors de la lecture des données de l'image :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des données de l'image"})
		return
	}
	photoVerso, _, err := image.Decode(&photoVersoData)
	if err != nil {
		log.Println("errer lors de image.Decode", err)
		return
	}
	outfile, err := os.Create("getPhotoVerso.png")
	if err != nil {
		log.Println("error lors de la creation de fichier", err)
		return
	}

	png.Encode(outfile, photoVerso)
	outfile.Close()

	var imgBuffer bytes.Buffer
	c.Header("Content-Type", "image/png")
	png.Encode(&imgBuffer, photoVerso)
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), imgBuffer.Bytes())

	log.Println("Image envoyée avec succès.")
}
