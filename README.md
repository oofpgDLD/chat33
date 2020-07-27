# chat33
## 快速开始
#### 环境准备
环境需求：linux-amd64

预留端口：8088
1. 安装go v1.13.0
```
# 教程
https://golang.org/doc/install?download=go1.14.6.linux-amd64.tar.gz
```
2. 安装docker
```
# 教程
https://docs.docker.com/engine/install/ubuntu/
```
3. 安装docker-compose
```
# 教程
https://docs.docker.com/compose/install/
```
#### Build&Run
```
cd $(GOPATH)/src/github.com/33cn/chat33 && make install
```
部署完毕后，控制台输入：`curl http://127.0.0.1:8088/chat/inner/ping`检测是否安装成功。

结果得到`success`表示部署成功
#### 卸载
```
cd $(GOPATH)/src/github.com/33cn/chat33 && make uninstall
```