FROM redis:4.0.2
USER root
COPY redis.conf /etc/redis/redis.conf
RUN mkdir -p -m 777 /var/log/redis/ /var/run/redis/ /var/lib/redis/
CMD [ "redis-server", "/etc/redis/redis.conf" ]