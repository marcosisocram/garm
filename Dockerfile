FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go mod download && go build -o /bin/main /app

FROM scratch

COPY --from=build /bin/main /

ENTRYPOINT ["/main"]
