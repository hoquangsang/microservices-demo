FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o /helloservice .

FROM gcr.io/distroless/static
COPY --from=builder /helloservice /helloservice
EXPOSE 8080
ENTRYPOINT ["/helloservice"]
