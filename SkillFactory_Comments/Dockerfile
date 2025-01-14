FROM golang:1.22 AS compiling_stage

WORKDIR /comments

COPY go.mod go.sum ./
RUN GODEBUG=http2client=0 go mod download
COPY . .

RUN go build -o comments ./main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=compiling_stage /comments .

COPY .env /root/.env

RUN chmod +x /root/comments

ENTRYPOINT ["./comments"]

EXPOSE 8882
