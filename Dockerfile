FROM golang:1.19

ENV CODECOV_TOKEN=7e965a4b-dc9c-4a3b-999a-f6b45f71cc84

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v ./...
RUN go test  -race -coverprofile=coverage.out -covermode=atomic

# Upload codecov-results
RUN curl -Os https://uploader.codecov.io/latest/linux/codecov && chmod +x codecov && ./codecov
