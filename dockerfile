FROM golang:latest
MAINTAINER "380562159@qq.com"
WORKDIR /go/src/azuremirror
ADD ./content /go/src/azuremirror/content
ADD ./main.go /go/src/azuremirror/
RUN go build .
EXPOSE 8010
ENTRYPOINT ["./azuremirror"]
