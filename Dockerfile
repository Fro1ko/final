FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /scheduler .

FROM alpine:latest

WORKDIR /app
COPY --from=build /scheduler ./scheduler
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

VOLUME ["/data"]
CMD ["./scheduler"]
