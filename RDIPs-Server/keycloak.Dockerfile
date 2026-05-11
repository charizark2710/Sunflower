FROM charizark2710/sunflower-migration as migration

FROM quay.io/keycloak/keycloak:26.6 as builder

# Enable health and metrics support
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

WORKDIR /opt/keycloak
RUN /opt/keycloak/bin/kc.sh build

# --- Optimized Runtime Stage ---
# Using Eclipse Temurin JRE on Alpine for the smallest, most secure footprint
FROM eclipse-temurin:21-jre-alpine

# Install dependencies needed for Keycloak and health checks
# Note: Alpine uses 'apk' instead of 'apt-get'
RUN apk add --no-cache curl bash

COPY --from=builder /opt/keycloak/ /opt/keycloak/
COPY --from=migration /migration/keycloak/ /opt/keycloak/data/import/

WORKDIR /opt/keycloak

# Security best practice: Create a dedicated group and user
RUN addgroup -g 1000 -S keycloak && \
    adduser -u 1000 -S keycloak -G keycloak && \
    chown -R keycloak:keycloak /opt/keycloak && \
    chmod -R 755 /opt/keycloak

USER keycloak

# Runtime Configurations
ENV KEYCLOAK_IMPORT=/opt/keycloak/data/import/RDIPs-realm.json
ENV KC_FEATURES=token-exchange
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]