FROM golang:1.15.2-alpine3.12 as build
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go build -o web .

FROM alpine:latest
WORKDIR /app
COPY --from=build /build/web .
EXPOSE 5000
CMD ["/app/web"]

