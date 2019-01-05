package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/gomodule/redigo/redis"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/srvc/fail"
)

var redisPool *redis.Pool

func init() {
	honeybadger.Configure(honeybadger.Configuration{APIKey: os.Getenv("HONEYBADGER_API_KEY")})

	redisPool = &redis.Pool{
		Wait: true,
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", "redis:6379") },
	}
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
		return err
	}

	subName := os.Getenv("PUBSUB_SUBSCRIPTION")
	if subName == "" {
		log.Fatal("Prease set `PUBSUB_SUBSCRIPTION` in ENV")
		return fmt.Errorf("Prease set `PUBSUB_SUBSCRIPTION` in ENV")
	}
	sub := client.Subscription(subName)

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		conn := redisPool.Get()
		defer conn.Close()

		id := string(m.Data)
		log.Printf("message received. id: %v", id)

		conn.Send("INCR", id)
		conn.Send("EXPIRE", id, 60*60)
		if err := conn.Flush(); err != nil {
			log.Printf("Failed to exec redis cmd incr & expire")
			honeybadger.Notify(err)
			m.Nack()
		}

		m.Ack()

		count, err := redis.Int64(conn.Do("GET", id))
		if err != nil {
			honeybadger.Notify(err, honeybadger.Context{"id": id, "subscription": subName})
			log.Printf("Failed to exec redis cmd get")
			return
		}

		if count > 1 {
			honeybadger.Notify(fail.New("Find duplicate"), honeybadger.Context{"id": id, "subscription": subName, "count": count})
			log.Printf("Find duplicate")
		}
	})

	if err != nil {
		return fail.Wrap(err)
	}

	return nil
}
