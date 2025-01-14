FROM golang:1.22 AS compiling_stage

WORKDIR /news

COPY go.mod go.sum ./
RUN GODEBUG=http2client=0 go mod download
COPY . .

RUN go build -o news ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=compiling_stage /news/news /root/news

COPY .env /root/.env
COPY config.json /root/config.json 

RUN chmod +x /root/news

ENTRYPOINT ["/root/news"]

EXPOSE 8883
