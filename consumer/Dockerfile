FROM golang:1.24 as builder
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer .

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app

COPY --from=builder /app/consumer .
COPY --from=builder /app/protogen ./protogen

# FOR QUICKTESTING PURPOSES ONLY, REMOVE IN PRODUCTION
# COPY --from=builder /app/config.yaml .

EXPOSE 9080

ENTRYPOINT ["/app/consumer"]