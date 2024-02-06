FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./goapp ./


FROM alpine:3
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --from=build /app/goapp ./
ENTRYPOINT ["./goapp"]
