FROM golang:1.13.4 as builder

WORKDIR /webhook
COPY . .
ENV CGO_ENABLED=0
RUN go build -o webhook

FROM scratch

COPY --from=builder /webhook/webhook /

CMD ["/webhook"]