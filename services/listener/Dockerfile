FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /listener ./services/listener/cmd

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build /listener /listener
EXPOSE 8090
ENTRYPOINT ["/listener"]
