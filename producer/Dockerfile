FROM golang:1.24 as builder
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o producer .

FROM gcr.io/distroless/static-debian12:debug
WORKDIR /app

COPY --from=builder /app/producer .
COPY --from=builder /app/protogen ./protogen

# FOR QUICKTESTING PURPOSES ONLY, REMOVE IN PRODUCTION
# COPY --from=builder /app/config.yaml .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/producer"]