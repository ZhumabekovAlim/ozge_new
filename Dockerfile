FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/web

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/app ./app
COPY --from=build /app/config ./config
COPY --from=build /app/static ./static
COPY --from=build /app/migrations ./migrations
COPY --from=build /app/DejaVuSans.ttf ./DejaVuSans.ttf
ENV GIN_MODE=release
EXPOSE 4000
CMD ["./app"]