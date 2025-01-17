FROM golang:1.23

COPY ./ /home/app

WORKDIR /home/app

RUN go build -o server -buildvcs=false

CMD /home/app/server