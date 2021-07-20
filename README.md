# Virtual Target

Virtual-Target（虚拟靶机）可以通过拖入一个响应文件，自动开一个新端口作为web服务端返回该响应

例如新建文件：modules/IIS.RESPONSE，直接启动`./virtual-target`会默认从1024端口开始启动服务，返回下方的响应

注意：无需设置Content-Length头，会自动计算并添加

```http request
HTTP/1.1 200 OK
Connection: close
Cache-Control: private
Content-Type: text/html
Date: Tue, 06 Jul 2021 03:48:58 GMT
Server: Microsoft-IIS/6.0
X-Powered-By: ASP.NET

<h1>IIS</h1>

```

如果1024端口被占用，将自动尝试下一个。可同时配置多个响应，效果如下

```shell
./virtual-target 
[info] [11:44:15] start DEFAULT at port 1024
[info] [11:44:15] start IIS at port 1025
[info] [11:44:15] start JQUERY at port 1026
[info] [11:44:15] start MEDUSA at port 1027
[info] [11:44:15] start NGINX at port 1028
[info] [11:44:15] start ORACLE at port 1029
[info] [11:44:15] start PHP at port 1030
[info] [11:44:15] start SUPERVISORD at port 1031
[info] [11:44:15] start TOMCAT at port 1032
[info] [11:44:15] start WORDPRESS at port 1033
```

也可以指定配置文件：`./virtual-target -c target.conf`

配置文件必须是`name=port`格式，其中name应与modules下的响应文件一致，总共的数量也必须一致
```text
DEFAULT=80
IIS=8001
JQUERY=8002
MEDUSA=8003
NGINX=8004
ORACLE=8005
PHP=8006
SUPERVISORD=8007
TOMCAT=8008
WORDPRESS=8009
```

使用Docker可以进行交叉编译：`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go`

注意指定host为0.0.0.0

```dockerfile
FROM centos:latest

ADD ./main /main
ADD ./target.conf /target.conf
ADD ./modules/ /modules/

EXPOSE 80
......
EXPOSE 8009

CMD /main -c /target.conf -H 0.0.0.0
```