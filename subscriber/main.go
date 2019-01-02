package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/srvc/fail"
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

	err := subscribe(ctx)

	if err != nil {
		log.Fatalln(err)
		honeybadger.Notify(err)
		return 1
	}

	return 0
}

func subscribe(ctx context.Context) error {
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	subName := os.Getenv("PUBSUB_SUBSCRIPTION")
	if subName == "" {
		log.Fatal("Prease set `PUBSUB_SUBSCRIPTION` in ENV")
	}
	sub := client.Subscription(subName)

	log.Println(sub)
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		// TODO Check uniq
		log.Printf("uuid: %s\n", m.Data)
		m.Ack()
	})

	if err != nil {
		return fail.Wrap(err)
	}

	return nil
}
