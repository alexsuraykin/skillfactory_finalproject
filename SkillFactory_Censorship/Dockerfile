FROM golang:1.22 AS compiling_stage

WORKDIR /censorship

COPY go.mod go.sum ./
RUN GODEBUG=http2client=0 go mod download
COPY . .

RUN go build -o censorship ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=compiling_stage /censorship/censorship /root/censorship

COPY .env /root/.env

RUN chmod +x /root/censorship

ENTRYPOINT ["/root/censorship"]

EXPOSE 8884
