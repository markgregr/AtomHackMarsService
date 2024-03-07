FROM golang:1.21.3-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o main

FROM alpine AS final

WORKDIR /

COPY --from=build-stage /app/cmd/main /main

ENTRYPOINT ["/main"]