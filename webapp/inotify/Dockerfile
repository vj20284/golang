FROM golang:alpine
ADD . /go/src/myapp
WORKDIR /go/src/myapp
RUN apk update && apk add git
RUN go get myapp
RUN go install
ENTRYPOINT ["/go/bin/myapp"]