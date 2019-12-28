FROM golang:1.13 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ARG GOOS=linux
ARG GIT_COMMIT=0
ARG GIT_VERSION=dev

WORKDIR /app
COPY . .

RUN go build -o /bin/nats-connector-example \
    -v -ldflags "-extldflags \"-static\"" \
    .

FROM alpine:3 as image

COPY --from=builder /bin/nats-connector-example /bin/nats-connector-example

CMD [ "/bin/nats-connector-example" ]
