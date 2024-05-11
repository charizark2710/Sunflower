FROM quay.io/keycloak/keycloak:latest as builder

# Enable health and metrics support
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

# Configure a database vendor
# ENV KC_DB=postgres

WORKDIR /opt/keycloak
RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:latest
COPY --from=builder /opt/keycloak/ /opt/keycloak/

USER root
RUN chown -R keycloak /opt
RUN chmod -R u+rwx /opt

USER keycloak
COPY ./RDIPs-Realm.json /opt/keycloak/data/import/
ENV KEYCLOAK_IMPORT=/opt/keycloak/data/import/RDIPs-realm.json
# change these values to point to a running postgres instance
# ENV KC_DB=postgres
# ENV KC_DB_URL=<DBURL>
# ENV KC_DB_USERNAME=<DBUSERNAME>
# ENV KC_DB_PASSWORD=<DBPASSWORD>
# ENV KC_HOSTNAME=localhost
# RUN /opt/keycloak/bin/kc.sh export --realm RDIPs --dir /opt --users realm_file
RUN ["/opt/keycloak/bin/kc.sh"]