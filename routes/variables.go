package routes

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"

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
var Bucket *gridfs.Bucket