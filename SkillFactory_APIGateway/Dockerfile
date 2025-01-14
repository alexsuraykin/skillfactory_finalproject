FROM golang:1.22 AS compiling_stage

WORKDIR /controller

COPY go.mod go.sum ./
RUN GODEBUG=http2client=0 go mod download
COPY . .

RUN go build -o controller ./main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=compiling_stage /controller .

COPY .env /root/.env

RUN chmod +x /root/controller

ENTRYPOINT ["./controller"]

EXPOSE 8080
