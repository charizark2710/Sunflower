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

USER rabbitmq

COPY --from=migration --chown=rabbitmq:rdips /migration/rabbitmq/rabbitmq.conf /etc/rabbitmq/rabbitmq.conf
COPY --from=migration --chown=rabbitmq:rdips /migration/rabbitmq/policy_definitions.json /rabbitmq/definitions/policy_definitions.json

RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_management && rabbitmq-plugins enable rabbitmq_mqtt