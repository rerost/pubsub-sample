package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/honeybadger-io/honeybadger-go"
)

func init() {
	honeybadger.Configure(honeybadger.Configuration{APIKey: os.Getenv("HONEYBADGER_API_KEY")})
}

func main() {
	statusCode := run()

	os.Exit(statusCode)
}

func run() int {
	defer honeybadger.Monitor()
	defer honeybadger.Flush()

	ctx := context.Background()

	err := sendManyMessage(ctx)

	if err != nil {
		honeybadger.Notify(err)
		return 1
	}

	return 0
}

func sendManyMessage(ctx context.Context) error {
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	topicName := os.Getenv("PUBSUB_TOPIC")
	if topicName == "" {
		log.Fatal("Prease set `PUBSUB_TOPIC` in ENV")
	}

	topic := client.Topic(topicName)
	defer topic.Stop()

	var results []*pubsub.PublishResult

	for i := 0; i < 10000; i++ {
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		r := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(u.String()),
		})
		results = append(results, r)
	}

	for _, r := range results {
		_, err := r.Get(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
