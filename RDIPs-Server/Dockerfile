FROM nginx
COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d /opt/nginx/conf.d
COPY stream.conf.d /opt/nginx/stream.conf.d

# RUN [ "/usr/sbin/nginx", "-c",  "/opt/nginx/conf"]
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
# COPY nginx.conf /etc/nginx/nginx.conf
# COPY ./html/ /usr/share/nginx/html/
