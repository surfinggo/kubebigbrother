FROM golang:1.16.4 as go

COPY . /working

WORKDIR /working

RUN go mod download

RUN go build -o /build/kbb .

FROM ubuntu:20.04

COPY --from=go /build/kbb /kbb

CMD /kbb