FROM rabbitmq:latest as base

WORKDIR /rabbitmq


RUN groupadd -f rdips
RUN chown -R rabbitmq:rdips /rabbitmq
RUN chmod -R gu+rwx /rabbitmq
RUN chown -R rabbitmq:rabbitmq /var/lib/rabbitmq

RUN mkdir -p /data/log
RUN mkdir -p /data/mnesia

RUN chown -R rabbitmq:rabbitmq /data/mnesia
RUN chmod -R gu+rwx /data/mnesia

RUN chown -R rabbitmq:rabbitmq /data/log
RUN chmod -R gu+rwx /data/log

USER rabbitmq

# Copy rabbitmq.conf
COPY --chown=rabbitmq:rdips rabbitmq.conf /etc/rabbitmq/rabbitmq.conf

RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_management