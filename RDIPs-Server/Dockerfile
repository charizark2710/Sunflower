FROM sunflower-ui as builder
WORKDIR /sunflower
CMD [ "serve", "-s", "build" ]

FROM nginx:stable-alpine
COPY --from=builder /sunflower/build /usr/share/nginx/html/sunflower/
COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d /opt/nginx/conf.d
COPY stream.conf.d /opt/nginx/stream.conf.d

# RUN [ "/usr/sbin/nginx", "-c",  "/opt/nginx/conf"]
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
# COPY nginx.conf /etc/nginx/nginx.conf
# COPY ./html/ /usr/share/nginx/html/
