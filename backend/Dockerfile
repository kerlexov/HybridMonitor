FROM golang:alpine as builder
ENV GO111MODULE=on
ENV GOPRIVATE=github.com/kerlexov/HybridMonitor/backend/
LABEL maintainer="<kerlexov@gmail.com>"

RUN apk update && apk add --no-cache git && mkdir /go/src/backend

WORKDIR /go/src/backend

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/backend/main .

EXPOSE 9393

CMD ["./main"]
