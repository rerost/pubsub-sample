package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalln(err)
	}

	topic := client.Topic(os.Getenv("PUBSUB_TOPIC"))
	if topic == nil {
		topic, err = client.CreateTopic(ctx, os.Getenv("PUBSUB_TOPIC"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	sub := client.Subscription(os.Getenv("PUBSUB_SUBSCRIPTION"))
	if sub == nil {
		_, err = client.CreateSubscription(
			ctx,
			os.Getenv("PUBSUB_SUBSCRIPTION"),
			pubsub.SubscriptionConfig{Topic: topic},
		)

		if err != nil {
			log.Fatal(err)
		}
	}
}
