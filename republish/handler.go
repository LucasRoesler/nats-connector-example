package function

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	nats "github.com/nats-io/nats.go"
	handler "github.com/openfaas-incubator/go-function-sdk"
)

// Handle a serverless request
func Handle(req handler.Request) (handler.Response, error) {
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
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not connect to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}
	defer nc.Close()

	msg := bytes.TrimSpace(req.Body)

	log.Printf("Publishing \"%s\" to: %s", string(msg), subject)
	err = nc.Publish(subject, msg)
	if err != nil {
		log.Println(err)
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not publish to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}

	return handler.Response{
		Body:       []byte(fmt.Sprintf("The msg: \"%s\" has been republished to \"%s\"", string(msg), subject)),
		StatusCode: http.StatusOK,
	}, nil
}
