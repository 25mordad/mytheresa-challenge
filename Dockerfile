FROM golang AS build

WORKDIR /go/src/app
COPY . .
COPY . /go/src/github.com/25mordad/mytheresa-challenge
RUN go get -d -v ./...


RUN go clean -testcache
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mytheresa-challenge cmd/server/main.go


FROM alpine:latest as certs
RUN apk --update add ca-certificates


FROM alpine:latest
EXPOSE 8081
RUN apk --update add tzdata
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /go/src/app/mytheresa-challenge .

CMD [ "/mytheresa-challenge" ]

#build image:
#docker build -t mytheresa-challenge -f Dockerfile .
#run image:
#docker run mytheresa-challenge
