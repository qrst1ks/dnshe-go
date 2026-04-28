FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w -X main.version=$VERSION" -o /out/dnshe-go .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates && adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=build /out/dnshe-go /usr/local/bin/dnshe-go
RUN mkdir -p /app/data && chown -R app:app /app
USER app
EXPOSE 9876
ENTRYPOINT ["dnshe-go"]
CMD ["-l", ":9876", "-c", "/app/data/config.json"]
