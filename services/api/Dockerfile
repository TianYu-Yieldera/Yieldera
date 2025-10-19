FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /api ./services/api/cmd

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates wget && update-ca-certificates
COPY --from=build /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
