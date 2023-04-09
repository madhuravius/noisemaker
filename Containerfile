FROM docker.io/golang:latest as build

WORKDIR /app
COPY . .
RUN make build

FROM scratch
COPY --from=build /app/bin/trashbin ./trashbin

ENTRYPOINT ["./trashbin"]