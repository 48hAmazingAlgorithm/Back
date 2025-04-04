package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	routes "challenge/routes"
)

func main() {
	routes.ConnectMongoDB()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/getClient", routes.GetClient)
	router.POST("/postClient", routes.PostClient)
	router.GET("/getIndividus", routes.GetIndividu)
	router.POST("/postIndividu", routes.PostIndividu)
	router.GET("/getPhotoRecto/:client_id", routes.GetPhotoRecto)
	router.GET("/getPhotoVerso/:client_id", routes.GetPhotoVerso)
	router.POST("/postPhotoRecto", routes.PostPhotoRecto)
	router.POST("/postPhotoVerso", routes.PostPhotoVerso)

	router.Run(":8080")
}
