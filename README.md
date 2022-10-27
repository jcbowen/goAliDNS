# goAliDNS

go 语言自动获取ip地址并更新到阿里动态DDNS里，实现动态解析

## 注意

* 初次运行的时候，如果没有配置文件，会自动生成配置文件，然后退出程序，需要手动修改配置文件

***

## 使用方法

Linux/Mac使用方法(Win差不多，只是需要的文件是.exe那个)：

```shell
./aliDDNS_linux
```

如果您想要后台运行，可以使用nohup命令

```shell
nohup ./aliDDNS_linux --log &
```

如果您想要将日志输出到文件而不是输出到控制台，可以传递参数--log

```shell
./aliDDNS_linux --log
```

***
目录说明

```
├─data // 配置/日志文件目录
│  │
│  └─conf.json // 配置信息
│
└─updateIp_l_linux //go编译出的可执行文件
```

***
conf.json结构如下，使用时请去掉注释

```json
{
 "AliOpenApiStruct": {
  "accessKeyId": "阿里云AccessKey ID",
  "accessKeySecret": "阿里云AccessKey Secret"
 },
 "subDomain": "www.example.com", // 需要解析的子域名
 "type": "A" // ip类型 A:ipv4, AAAA:ipv6, ALL:ipv4和ipv6
}
```