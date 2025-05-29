package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBinstance initialise une connexion MongoDB et retourne un client
func DBinstance() *mongo.Client {
	// Charger les variables d'environnement (optionnel)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Récupérer l'URL MongoDB
	MongoDb := os.Getenv("MONGODB_URL")
	if MongoDb == "" {
		// Fallback vers MONGODB_URI si MONGODB_URL n'est pas défini
		MongoDb = os.Getenv("MONGODB_URI")
		if MongoDb == "" {
			log.Fatal("Neither MONGODB_URL nor MONGODB_URI is set in environment variables")
		}
	}

	// Créer un nouveau client MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Contexte avec un timeout pour la connexion
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connecter le client à MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

// Client est une instance globale de MongoDB
var Client *mongo.Client = DBinstance()

// OpenCollection retourne une collection MongoDB
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME") // Récupérer le nom de la base de données
	if dbName == "" {
		log.Fatal("DB_NAME is not set in environment variables")
	}

	// Accéder à la collection
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}
