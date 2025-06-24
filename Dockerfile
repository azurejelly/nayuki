FROM golang:1.24.4 AS build

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=build /build/nayuki .

ENTRYPOINT [ "./nayuki" ]