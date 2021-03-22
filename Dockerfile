FROM golang:1.12-alpine as builder
WORKDIR /app

COPY . .
RUN go build -o bin/

FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/img.farm /usr/local/bin/
CMD ["img.farm"]