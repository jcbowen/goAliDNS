# goAliDNS
go 语言自动获取ip地址并更新到阿里动态DNS里，实现动态解析

###执行文件同目录下需要创建一个config目录，并在其中放一个config.json

***
目录说明
```
├─config // 配置文件目录
│  │
│  └─config.json // 配置信息
│
└─updateIp_l_linux //go编译出的可执行文件
```

***
config.json结构如下
```
{
  "aliOpenApi": {
    "accessKeyId": "阿里云ACCESSKEYID",
    "accessKeySecret": "阿里云ACCESSKEYSECRET"
  },
  "subDomain": "www.domain.com", // 需要修改的二级域名
  "setting": {
    "type": "6" // 解析类型，请参考阿里云动态解析文档
  }
}
```