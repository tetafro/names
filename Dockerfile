FROM golang:1.18-alpine3.16 AS build

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev

COPY . .

RUN go build -o ./bin/names .

FROM alpine:3.16

WORKDIR /app

COPY --from=build /build/bin/names /app/
COPY index.html /app/index.html
COPY static /app/static

RUN apk add --no-cache ca-certificates && \
    addgroup -S -g 5000 names && \
    adduser -S -u 5000 -G names names && \
    chown -R names:names .

USER names

EXPOSE 8080

ENTRYPOINT ["/app/names", "-port", "8080"]
