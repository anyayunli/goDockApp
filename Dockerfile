FROM golang:latest
RUN mkdir /go/src/goDockApp
ADD . /go/src/goDockApp
WORKDIR /go/src/goDockApp
RUN ls -l
RUN go get ./...
RUN go build -o main .
EXPOSE 3344
CMD ["/go/src/goDockApp/main"]
