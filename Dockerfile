# docker run --rm -it --expose 8003 -p 8003:8003 -p 8004:8004 -t kiki
FROM golang:alpine
RUN apk add --no-cache git g++

MAINTAINER Zack Scholl "zack.scholl@gmail.com"

RUN mkdir /app
WORKDIR /app

RUN go get -v github.com/schollz/kiki
RUN go install -v github.com/schollz/kiki
RUN mv /go/bin/kiki /app/kiki
RUN rm -rf /go
RUN apk del git g++
ENTRYPOINT ["/app/kiki","-path","/app","-no-browser","-expose"]
