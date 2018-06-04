FROM golang:1.10 as build

ADD . /go/src/github.com/DazWilkin
WORKDIR /go/src/github.com/DazWilkin

RUN go get ./hellohenry
RUN go get ./healthcheck

FROM gcr.io/distroless/base as runtime

LABEL maintainer="Daz Wilkin"

COPY --from=build /go/bin/hellohenry /
COPY --from=build /go/bin/healthcheck /

ENTRYPOINT ["/hellohenry"]

HEALTHCHECK --interval=5s --timeout=5s CMD ["/healthcheck","http://localhost:8080/healthz"]