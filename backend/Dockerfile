FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy module definition and source, then resolve deps (creates go.sum in the image)
COPY go.mod ./
COPY . .
RUN go mod tidy && go build -o server ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]