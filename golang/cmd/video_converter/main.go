package main

import (
	"database/sql"
	"fmt"
	"go-converter/internal/converter"
	"go-converter/internal/rabbitmq"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func main() {
	db, err := connectPostgres()
	if err != nil {
		panic(err)
	}

	rabbitmqUrl := getEnvOrDefault("RABBITMQ_URL", "amqp://guest:gust@rabbitmq:5672/")
	rabbitClient, err := rabbitmq.NewRabbitClient(rabbitmqUrl)
	if err != nil {
		panic(err)
	}
	defer rabbitClient.Close()

	conversionExchange := getEnvOrDefault("CONVERSION_EXCHANGE", "conversion_exchange")
	queueName := getEnvOrDefault("CONVERSION_QUEUE", "video_conversion_queue")
	conversionKey := getEnvOrDefault("CONVERSION_KEY", "conversion")
	confirmationQueue := getEnvOrDefault("CONFIRMATION_QUEUE", "video_confirmation_queue")
	confirmationKey := getEnvOrDefault("CONFIRMATION_KEY", "finish-conversion")

	vc := converter.NewVideoConverter(rabbitClient, db)
	// vc.Handle([]byte(`{"video_id": 2, "path": "/media/uploads/2"}`))

	messages, err := rabbitClient.ConsumeMessages(conversionExchange, conversionKey, queueName)
	if err != nil {
		// This not stop reading other messages
		slog.Error("Failed to consume messages", slog.String("error", err.Error()))
	}

	slog.Info("[*] Waiting for messages")

	for message := range messages {
		go func(delivery amqp.Delivery) {
			vc.Handle(delivery, conversionExchange, confirmationKey, confirmationQueue)
		}(message)
	}
}

func connectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "user")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "password")
	dbName := getEnvOrDefault("POSTGRES_DB", "converter")
	host := getEnvOrDefault("POSTGRES_HOST", "postgres")
	sslMode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", user, password, dbName, host, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Error connecting to database", slog.String("connStr", connStr))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Error pinging database", slog.String("connStr", connStr))
		return nil, err
	}

	slog.Info("Connected to Postgres successfully")
	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
