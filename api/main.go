package main

import (
	"errors"
	"github.com/Kamva/mgm"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/toniastro/h_ai/api/databases/redis"
	"github.com/toniastro/h_ai/api/handlers"
	"github.com/toniastro/h_ai/api/queues/rabbitMq"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("PORT not set in .env")
	}

	host, exist := os.LookupEnv("HOST")
	if !exist {
		//log.Fatal("HOST not set in .env")
	}
	_, route := setUp()
	srv := &http.Server{
		Addr:         host + ":" + port,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      route,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func setUp() (*handlers.Handler, *mux.Router) {
	connectionString := os.Getenv("DATABASE_URL")

	err := mgm.SetDefaultConfig(nil, "hasty_ai", options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Println(errors.New("cannot set mgm config"))
		log.Fatal(err)
	}

	redisAddress, exist := os.LookupEnv("REDIS_ADDRESS")
	if !exist {
		log.Fatal("REDIS_ADDRESS not set in .env")
	}

	redisPassword, exist := os.LookupEnv("REDIS_PASSWORD")
	if !exist {
		log.Println("REDIS_PASSWORD not set in .env")
	}

	//Initiate Redis
	redis := redis.New(&redis.Config{
		Addr:     redisAddress,
		Password: redisPassword,
	})

	rabbitHost, exist := os.LookupEnv("RABBIT_URL")
	if !exist {
		log.Fatal("RABBIT_URL not set in .env")
	}

	queueName, exist := os.LookupEnv("RABBIT_QUEUE_NAME")
	if !exist {
		log.Fatal("RABBIT_QUEUE_NAME not set in .env")
	}

	//Initiate RabbitMQ
	rabbit := rabbitMq.New(rabbitHost, queueName)

	//Initiate Routes
	route := mux.NewRouter()
	h := handlers.New(redis, rabbit)
	h.AddRoutes(route)
	return h, route
}
