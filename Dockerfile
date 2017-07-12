FROM golang:latest

WORKDIR /go/src/app
COPY . .

# "gogo"
RUN go-wrapper download  
# "go install -v ./..."
#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
ENV CGO_ENABLED 0
ENV GOOS linux
RUN go-wrapper install -a -installsuffix cgo 

FROM scratch
COPY --from=0 /go/bin/app .
CMD ["/app"]