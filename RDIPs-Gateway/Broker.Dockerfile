FROM rabbitmq:latest as base

WORKDIR /rabbitmq

RUN groupadd -f rdips
RUN chown -R rabbitmq:rdips /rabbitmq
RUN chmod -R gu+rwx /rabbitmq
USER rabbitmq

# Copy rabbitmq.conf
COPY --chown=rabbitmq:rdips rabbitmq.conf /etc/rabbitmq/rabbitmq.conf
RUN mkdir -p /var/lib/rabbitmq/mnesia
RUN mkdir -p /var/log/rabbitmq
RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_management