FROM golang

WORKDIR /go/src/app
COPY . .

VOLUME ["/opt/go/QA/static", "/opt/go/QA/xlsx"]

EXPOSE 8080
CMD ["./main"]