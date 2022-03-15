FROM golang:1.17.8-alpine3.15

RUN mkdir endtermGolang

WORKDIR /endtermGolang

COPY . .

RUN export GO111MODULE=on
RUN cd /endtermGolang
RUN go get github.com/gin-gonic/gin
RUN go get golang.org/x/crypto/bcrypt
RUN go build -o main.exe

EXPOSE 8080

CMD [ "/endtermGolang/main.exe" ]

