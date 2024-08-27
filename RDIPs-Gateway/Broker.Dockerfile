FROM charizark2710/sunflower-migration as migration

FROM rabbitmq:3.13.5-management-alpine as base

WORKDIR /rabbitmq

RUN addgroup rdips

RUN chown -R rabbitmq:rdips /rabbitmq
RUN chown -R rabbitmq:rdips /var/lib/rabbitmq

RUN mkdir -p /data/log
RUN mkdir -p /data/mnesia

RUN chown -R rabbitmq:rdips /data/mnesia
RUN chmod -R gu+rwx /data/mnesia

RUN chown -R rabbitmq:rdips /data/log
RUN chmod -R gu+rwx /data/log

RUN apk update && apk add envsubst

USER rabbitmq
ARG BROKER_USER
ARG BROKER_PASSWORD

ENV BROKER_USER=${BROKER_USER}
ENV BROKER_PASSWORD=${BROKER_PASSWORD}

COPY --from=migration --chown=rabbitmq:rdips /migration/rabbitmq/rabbitmq.conf /etc/rabbitmq/rabbitmq-tmp.conf
RUN envsubst "$(printf '${%s} ' $(env | cut -d'=' -f1))" < /etc/rabbitmq/rabbitmq-tmp.conf > /etc/rabbitmq/rabbitmq.conf

RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_management && rabbitmq-plugins enable rabbitmq_mqtt
