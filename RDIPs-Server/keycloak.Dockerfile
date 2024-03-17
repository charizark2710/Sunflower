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

COPY ./RDIPs-Realm.json /opt/keycloak/imports
ENV KEYCLOAK_IMPORT=/opt/keycloak/imports/RDIPs-Realm.json
# change these values to point to a running postgres instance
# ENV KC_DB=postgres
# ENV KC_DB_URL=<DBURL>
# ENV KC_DB_USERNAME=<DBUSERNAME>
# ENV KC_DB_PASSWORD=<DBPASSWORD>
# ENV KC_HOSTNAME=localhost
RUN ["/opt/keycloak/bin/kc.sh"]