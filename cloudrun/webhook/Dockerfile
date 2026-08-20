FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./goapp ./main.go

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /app/goapp /
ENTRYPOINT ["/goapp"]
