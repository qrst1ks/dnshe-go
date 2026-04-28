FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/dnshe-go .

FROM alpine:3.22

RUN adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=build /out/dnshe-go /usr/local/bin/dnshe-go
RUN mkdir -p /app/data && chown -R app:app /app
USER app
EXPOSE 9876
ENTRYPOINT ["dnshe-go"]
CMD ["-l", ":9876", "-c", "/app/data/config.json"]

