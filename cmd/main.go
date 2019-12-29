package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	go func() {
		<-killSignal
		log.Println("Stopping...")
		cancel()
	}()

	subject := "faas-req"
	val, ok := os.LookupEnv("faas_request_subject")
	if ok {
		subject = val
	}

	respSubject := "faas-resp"
	val, ok = os.LookupEnv("faas_response_subject")
	if ok {
		respSubject = val
	}

	msg := "Hello World"
	val, ok = os.LookupEnv("faas_msg")
	if ok {
		subject = val
	}

	natsURL := nats.DefaultURL
	val, ok = os.LookupEnv("nats_url")
	if ok {
		natsURL = val
	}

	until := 60 * time.Second
	val, ok = os.LookupEnv("msg_until")
	if ok {
		until, err = time.ParseDuration(val)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Will send a message \"%s\" every 1s for %s\n", msg, until)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	go sendMessages(ctx, nc, subject, msg, until)

	err = listen(ctx, nc, respSubject, until)
	if err != nil {
		log.Fatalf("listening error: %s", err)
	}

	log.Println("Success!")
}

func listen(ctx context.Context, nc *nats.Conn, subject string, until time.Duration) error {
	var count int
	expectedCount := int(until.Seconds())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second+until)
	defer cancel()

	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	fmt.Printf("Listener started, expecting %d messages\n", expectedCount)
	for count = 1; count <= expectedCount; count++ {
		if ctx.Err() != nil {
			return err
		}
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Response %d of %d: \"%s\"\n", count, expectedCount, string(msg.Data))
	}

	return nil
}

func sendMessages(ctx context.Context, nc *nats.Conn, subject, msg string, until time.Duration) {
	t := time.NewTicker(time.Second)
	ctx, cancel := context.WithTimeout(ctx, until)
	defer cancel()

	rawMsg := []byte(msg)

	for {
		select {
		case <-ctx.Done():
			log.Println("Finished.")
			return
		case <-t.C:
			log.Printf("Sending \"%s\" to \"%s\"\n", msg, subject)
			err := nc.Publish(subject, rawMsg)
			if err != nil {
				log.Print(err)
				return
			}
		}
	}
}
