FROM golang:1.21.3 as base

WORKDIR /sunflower

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN useradd -ms /bin/bash sunflower
RUN groupadd -f rdips
RUN chown -R sunflower:rdips /sunflower
RUN chmod -R gu+rwx /sunflower

EXPOSE 8080

FROM base as production
COPY --chown=sunflower:rdips . .
USER sunflower
RUN go build -o sunflower
CMD [ "./sunflower" ]
