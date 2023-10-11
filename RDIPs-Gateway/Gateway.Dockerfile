FROM gcc:13.2.0

WORKDIR /RDIPs-Gateway

RUN useradd RDIPs-Gateway
RUN groupadd -f rdips
RUN chown -R RDIPs-Gateway:rdips /RDIPs-Gateway
RUN chmod -R gu+rwx /RDIPs-Gateway
USER RDIPs-Gateway

# Copy broker cpp
COPY --chown=RDIPs-Gateway:rdips . .
RUN ./premake5 gmake2
RUN make -C ./build
# ENTRYPOINT ["./bin/RDIPs-Gateway"]
