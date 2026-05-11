FROM golang:1.24.0 AS build

WORKDIR /dashboard

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rdips-dashboard .

FROM alpine:3.22.4 AS production

RUN addgroup -g 1000 rdips && \
    adduser -u 1000 -G rdips -D sunflower

WORKDIR /dashboard

COPY --from=build --chown=sunflower:rdips /dashboard/rdips-dashboard ./rdips-dashboard
RUN chmod +x rdips-dashboard

USER sunflower

EXPOSE 3000
ENV PORT=3000

CMD ["./rdips-dashboard"]
