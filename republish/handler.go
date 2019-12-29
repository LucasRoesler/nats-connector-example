package function

import (
	"fmt"
	"log"
	"os"

	nats "github.com/nats-io/nats.go"
)

// Handle a serverless request
func Handle(req []byte) string {
	subject := "faas-response"
	val, ok := os.LookupEnv("target_subject")
	if ok {
		subject = val
	}

	natsURL := nats.DefaultURL
	val, ok = os.LookupEnv("nats_url")
	if ok {
		natsURL = val
	}

	log.Printf("Connecting to: %s", natsURL)
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer nc.Close()

	log.Printf("Publishing \"%s\" to: %s", string(req), subject)
	err = nc.Publish(subject, req)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return fmt.Sprintf("The msg: \"%s\" has been republished to \"%s\"", string(req), subject)
}
