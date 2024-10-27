FROM golang:alpine AS build
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./go-pocket-link ./cmd/app/main.go

FROM alpine AS release
WORKDIR /app

RUN --mount=type=cache,target=/var/cache/apk apk --update \
  add ca-certificates tzdata && update-ca-certificates

EXPOSE 8080

COPY ./config/dev.yml ./config.yml
COPY --from=build /app/go-pocket-link .

ENTRYPOINT [ "./go-pocket-link" ]
