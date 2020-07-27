FROM golang:1.13
WORKDIR /app
COPY build/ /app
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct go build -v -o chat33

FROM scratch as final
COPY --from=build /app/build/chat33 .
CMD [ "./chat33" ]