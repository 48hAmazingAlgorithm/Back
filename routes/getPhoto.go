package routes

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPhotoRecto(c *gin.Context) {
	clientId := c.Param("client_id")
	individu, _, _, _, cancel := GetSingleIndividu(clientId, c)
	defer cancel()

	log.Println("Individu trouvé :", individu)

	photoRectoId, err := DecryptID(individu.PhotoRectoID)
	if err != nil {
		log.Println("Erreur lors du décryptage de l'ID de la photoRecto :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du décryptage de l'ID de la photoRecto"})
		return
	}
	photoRecto := getPhoto(photoRectoId, c)
	var imgBuffer bytes.Buffer
	c.Header("Content-Type", "image/png")
	png.Encode(&imgBuffer, photoRecto)
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), imgBuffer.Bytes())

	log.Println("Image envoyée avec succès.")
}

func GetPhotoVerso(c *gin.Context) {
	clientId := c.Param("client_id")
	individu, _, _, _, cancel := GetSingleIndividu(clientId, c)
	defer cancel()

	log.Println("Individu trouvé :", individu)

	photoVersoId, err := DecryptID(individu.PhotoVersoID)
	if err != nil {
		log.Println("Erreur lors du décryptage de l'ID de la photoRecto :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du décryptage de l'ID de la photoRecto"})
		return
	}
	photoVerso := getPhoto(photoVersoId, c)

	var imgBuffer bytes.Buffer
	c.Header("Content-Type", "image/png")
	png.Encode(&imgBuffer, photoVerso)
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), imgBuffer.Bytes())

	log.Println("Image envoyée avec succès.")
}

func getPhoto(id string, c *gin.Context) image.Image {
	log.Println("photoId décrypté :", id)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Erreur lors de la conversion de l'ID de la photo en ObjectId :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de la photo invalide"})
		return nil
	}
	log.Println("id de la photo", objectId)

	downloadStream, err := Bucket.OpenDownloadStream(objectId)
	if err != nil {
		log.Println("Erreur lors du téléchargement de la photo depuis GridFS :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du téléchargement de la photo"})
		return nil
	}
	defer downloadStream.Close()

	var photoData bytes.Buffer
	_, err = photoData.ReadFrom(downloadStream)
	if err != nil {

		log.Println("Erreur lors de la lecture des données de l'image :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture des données de l'image"})
		return nil
	}
	photo, _, err := image.Decode(&photoData)
	if err != nil {
		log.Println("errer lors de image.Decode", err)
		return nil
	}
	return photo
}
