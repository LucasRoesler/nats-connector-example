package function

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	nats "github.com/nats-io/nats.go"
	handler "github.com/openfaas-incubator/go-function-sdk"
)

// Handle a serverless request
func Handle(req handler.Request) (handler.Response, error) {
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
		msg = val
	}

	if len(req.Body) > 0 {
		msg = string(bytes.TrimSpace(req.Body))
	}

	natsURL := nats.DefaultURL
	val, ok = os.LookupEnv("nats_url")
	if ok {
		natsURL = val
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not connect to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}
	defer nc.Close()

	log.Printf("Sending \"%s\" to \"%s\"\n", msg, subject)
	err = nc.Publish(subject, []byte(msg))
	if err != nil {
		log.Println(err)
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not publish to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}

	log.Println("Waiting for response")
	sub, err := nc.SubscribeSync(respSubject)
	if err != nil {
		log.Println(err)
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not subscribe to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}
	defer sub.Unsubscribe()

	resp, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		log.Println(err)
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("failed waiting for response: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}

	if string(resp.Data) != msg {
		log.Printf("expected %s, got %s", msg, string(resp.Data))
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("invalid response received: %s", string(resp.Data))),
			StatusCode: http.StatusConflict,
		}
		return r, err
	}

	return handler.Response{
		Body:       []byte("Success!"),
		StatusCode: http.StatusOK,
	}, nil
}
