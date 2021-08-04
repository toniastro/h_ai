package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kamva/mgm"
	"github.com/joho/godotenv"
	"github.com/toniastro/h_ai/consumer/queues/rabbitMq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Jobs struct {
	mgm.DefaultModel `bson:",inline"`
	Status           string `json:"status" bson:"status"`
	TimeTaken        int    `bson:"time_taken" json:"time_taken"`
	ObjectID         string `bson:"object_id" json:"object_id"`
	JobID            string `bson:"job_id" json:"job_id"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	connectionString := os.Getenv("DATABASE_URL")

	err := mgm.SetDefaultConfig(nil, "hasty_ai", options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Println(errors.New("cannot set mgm config"))
		log.Fatal(err)
	}

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

	err = ConsumeJob(context.Background(), rabbit)
	if err != nil {
		log.Fatal(err)
	}
}

func ConsumeJob(ctx context.Context, rabbit *rabbitMq.Rabbit) error {

	timeConfigured := 25
	messageChannel, err := rabbit.Channel.Consume(
		rabbit.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Could not register consumer")
		log.Fatal(err)
	}

	stopChan := make(chan bool)

	go func() {

		log.Printf("Consumer ready, PID: %d", os.Getpid())

		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			jobs := &Jobs{}

			err = json.Unmarshal(d.Body, jobs)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			err = updateJob(ctx, jobs, timeConfigured)
			if err != nil {
				log.Println(err)
			}

			if err = d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	// Stop for program termination
	<-stopChan
	return nil
}

func updateJob(ctx context.Context, jobs *Jobs, timeout int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	select {
	case <-time.After(time.Second * time.Duration(jobs.TimeTaken)):
		return mongoUpdate(jobs.JobID, "COMPLETED")
	case <-ctx.Done():
		log.Println("Timed Out")
		return mongoUpdate(jobs.JobID, "CANCELLED")
	}
}

func mongoUpdate(jobId, status string) error {

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"job_id", jobId}}
	update := bson.D{{"$set", bson.D{{"status", status}}}}

	_, err := mgm.Coll(&Jobs{}).UpdateOne(context.TODO(), filter, update, opts)

	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Updated task as completed")
	return nil
}
