FROM ubuntu:16.04
WORKDIR /app
RUN mkdir -p -m 777 /app/etc /app/bin
ADD build/ .
ADD etc/ /app/etc
EXPOSE 8090
ENTRYPOINT  ["./chat33"]