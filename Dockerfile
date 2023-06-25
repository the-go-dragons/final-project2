FROM golang:1.20 AS go-build

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=go-build /app .

# Add wait-for-it.sh script
# ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh .
COPY ./wait-for-it.sh .
RUN chmod +x wait-for-it.sh