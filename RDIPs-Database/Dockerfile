FROM charizark2710/sunflower-migration as migration
WORKDIR /migration

FROM postgres:alpine3.23

COPY ./post-init.sh /docker-entrypoint-initdb.d/
COPY --from=migration /migration/postgres/ /docker-entrypoint-initdb.d/
RUN ["chmod", "+x", "/docker-entrypoint-initdb.d/post-init.sh"]


EXPOSE ${POSTGRES_PORT}

USER postgres
CMD [ "postgres" ]
