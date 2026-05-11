FROM golang:1.24.0 AS build

WORKDIR /sunflower

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sunflower-bin .

FROM alpine:3.22.4 AS production
# RUN useradd -u 1000 -ms /bin/bash sunflower
# RUN groupadd -f rdips -g 1000
# RUN chown -R sunflower:rdips /sunflower
# RUN chmod -R gu+rwx /sunflower
RUN addgroup -g 1000 rdips && \
    adduser -u 1000 -G rdips -D sunflower

WORKDIR /sunflower

COPY --from=build --chown=sunflower:rdips /sunflower/sunflower-bin ./sunflower
RUN chmod +x sunflower

USER sunflower

EXPOSE 8080

CMD [ "./sunflower" ]
