FROM golang:latest
MAINTAINER James Vipond
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/app/MTA.go" ]