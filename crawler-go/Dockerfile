FROM golang:1.9.2

ENV GOPATH /crawler

ADD src /crawler
RUN mv /crawler/vendor /crawler/src

WORKDIR /crawler

CMD ["make", "start"]