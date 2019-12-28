package main

import (
	"context"
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
			log.Println(err)
			return
		}
	}

	log.Printf("Will send a message \"%s\" every 1s for %s\n", msg, until)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer nc.Close()

	sendMessages(ctx, nc, subject, msg, until)

	log.Println("Finished.")

	err = nc.Drain()
	if err != nil {
		log.Println(err)
	}
}

func sendMessages(ctx context.Context, nc *nats.Conn, subject, msg string, until time.Duration) {
	t := time.NewTicker(time.Second)
	ctx, cancel := context.WithTimeout(ctx, until)
	defer cancel()

	rawMsg := []byte(msg)

	for {
		select {
		case <-ctx.Done():
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
