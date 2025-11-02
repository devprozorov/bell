package config

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

var DB *mongo.Database

func ConnectDB() *mongo.Database {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Читаем переменные
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	if uri == "" {
		log.Fatal("❌ MONGO_URI is not set in .env")
	}
	if dbName == "" {
		log.Fatal("❌ MONGO_DB_NAME is not set in .env")
	}

	fmt.Println("🔗 Connecting to MongoDB:", uri)

	// Настраиваем клиента
	clientOptions := options.Client().ApplyURI(uri)

	// Подключаемся
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Mongo connect error:", err)
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Mongo ping error:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")

	DB = client.Database(dbName)
	return DB
}
