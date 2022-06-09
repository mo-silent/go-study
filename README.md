> Author mogd 2022-04-08
> \
> Update mogd 2022-06-09

# go-study
go语言学习-比较乱

## k8s-client-go

对 Kubernetes 集群的一些 CURD 操作练习

## test

遇到一些好的代码用来测试结果，请忽略这个

## benchmark
这个是一个基准压测案例

## go-grpc-helloworld

官方 gRPC 案例

## go-grpc-CA

基于 CA 的 TLS 证书认证的 gRPC 

## go-grpc-encryptChat

一个加密聊天室案例，通过 CA 证书加密 RPC 通信。不过显示界面还需要优化

## go-redis

go-redis 的使用

## go-redis-shake-decode

redis-shake decode 模式下，读取 rbd 文件的流程简化

## go-build-system 

go-build-system 指定 go 文件生成不同操作系统的可执行文件

## go-mongodb

go-mongodb go 操作 mongoDB 的案例，计划形成接口，让前端调用，能够直接在界面上操作 mongodb

## go-http-getIP

> 个人抓包到的 IP 来解析实际属于哪个云厂商，读者可忽略

go-http-getIP go 使用 net/http 包请求 [IP 查询网站](http://mip.chinaz.com/?query=)来正则拿到 IP 对应的物理地址

命令行传参 `go run main.go -infile=D:/文件/2022-05-12/tmp.txt -outfile=D:/文件/2022-05-12/tmp.csv`
1. `-infile` 要查询的 I P列表，`txt` 或者 `csv`，没做类型判断
2. `outfile` 要输出的信息文件路径，`csv` 格式，没做类型判断

## go-ping

go-ping go 使用 ping 包来获取域名对应的 IP 地址

命令行传参 `go run main.go -infile=D:/文件/2022-05-12/tmp.txt -outfile=D:/文件/2022-05-12/tmp.csv`
1. `-infile` 要查询的域名列表，`txt` 或者 `csv`，没做类型判断
2. `outfile` 要输出的信息文件路径，`csv` 格式，没做类型判断

## go-smokeping-conf
> 个人用途，忽略

go-smokeping-conf 读取 CSV 文件，生成 smokeping 的配置文件

## go-cloud-ping

> 个人用途，忽略

go-cloud-ping 读取文件，获取ping 延迟和丢包

参数说明：
1. `infile` 要测试的 IP 列表，一个 IP 一行就行，CSV 格式
2. `outpath` 结果输出的路径
3. `region` 结果输出文件的前缀
4. `log` 程序运行日志的输出
5. `core` CPU 内核数量，默认为 4
6. `channel` 协程数量，默认为 100

## go-zabbix-sender

go-zabbix-sender 使用 exec 执行 linux 命令，调用 zabbix_sender 发送历史数据到 zabbix，用于 zabbix 历史数据迁移

参数说明：
1. `-dir` 历史数据文件存储的目录
2. `log` 程序运行日志的输出

```shell
go run main.go -dir=/root/test/ -log=./go-zabbix-sender.log
```

## PAT 

PAT 浙大的PAT真题练习代码