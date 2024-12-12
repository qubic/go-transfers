FROM golang:1.23 AS builder
ENV CGO_ENABLED=0

WORKDIR /src/go-transfers
COPY . /src/go-transfers

RUN go mod tidy
WORKDIR /src/go-transfers
RUN go build

FROM ubuntu:latest
LABEL authors="mio@qubic.org"

# copy executable from build stage
COPY --from=builder /src/go-transfers/go-transfers /app/go-transfers
# copy default configuration
COPY .env /app/

RUN chmod +x /app/go-transfers

WORKDIR /app

ENTRYPOINT ["./go-transfers"]