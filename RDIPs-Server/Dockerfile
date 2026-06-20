FROM charizark2710/sunflower-ui as builder
WORKDIR /sunflower
CMD [ "serve", "-s", "build" ]

FROM nginx:stable-alpine AS module-builder

# Install compilation tools and dependencies
RUN apk add --no-cache \
    git \
    gcc \
    libc-dev \
    make \
    pcre-dev \
    zlib-dev \
    openssl-dev \
    linux-headers

# 1. Clone the memc module source code
RUN git clone https://github.com/openresty/memc-nginx-module /tmp/memc-nginx-module

# 2. Download the Nginx source code that EXACTLY matches the running version
RUN NGINX_VERSION=$(nginx -v 2>&1 | cut -d'/' -f2) && \
    wget http://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz -O /tmp/nginx.tar.gz && \
    tar -zxvf /tmp/nginx.tar.gz -C /tmp

# 3. Configure and compile ONLY the dynamic module using official binary arguments
RUN NGINX_VERSION=$(nginx -v 2>&1 | cut -d'/' -f2) && \
    cd /tmp/nginx-${NGINX_VERSION} && \
    # Extract original compilation arguments so our module is binary-compatible
    CONFARGS=$(nginx -V 2>&1 | grep 'configure arguments:' | sed 's/configure arguments://') && \
    eval ./configure $CONFARGS --add-dynamic-module=/tmp/memc-nginx-module && \
    make modules && \
    cp objs/ngx_http_memc_module.so /tmp/

FROM nginx:stable-alpine
COPY --from=builder /sunflower/build /usr/share/nginx/html/sunflower/
COPY --from=module-builder /tmp/ngx_http_memc_module.so /etc/nginx/modules/ngx_http_memc_module.so

COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d /opt/nginx/conf.d
COPY stream.conf.d /opt/nginx/stream.conf.d

# RUN [ "/usr/sbin/nginx", "-c",  "/opt/nginx/conf"]
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
# COPY nginx.conf /etc/nginx/nginx.conf
# COPY ./html/ /usr/share/nginx/html/
