## 介绍

在某些场景下，需要访问指定页面或者资源前做一些逻辑判断，不想破话原有的业务逻辑或者功能，只想在访问的链路中加一个可以用API的方式来判断。

## 安装

> 源码编译安装

```sh
git clone xxxxxx
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build caddy.go
```

> 直接下载程序安装

## 使用

配置文件：

```
{
  http_port     9080
  https_port    9443
}

:9080 {
  route {
    logic_api disabled=no tag="jira" url="http://172.26.5.164:8888/api/authSite/checkAuthSite" redirect="http://www.baidu.com"
    reverse_proxy 172.26.1.51:8080
  }
}

提示：
url 是一个判断的接口：
接口入参：
{
	"address": "", //对应的是IP+端口
	"cookies": [   //对应的网站回传的cookies的值
		{
			"name": "",
			"value": ""
		}
	],
	"host": "",  //对应的是被访问的Host地址
	"proto": "", //对应的是协议
	"agent": "", //对应浏览器的user-agent的值
	"tag": ""  //对应的配置文件的tag的值，可以对不同的route的配置不同的标识，不同场景可以根据标识来做过滤
}
接口响应：这个是固定的格式，需要在写远程验证接口的时候构建
{
  "code": 200,
  "success": true,  // 这个是主要来判断远程验证的结果，true为通过，false为调整页面
  "message": "操作成功",
  "data": null,
  "time": "2021-01-15 11:44:18"
}

redirect 是失败后的一个跳转地址：

```



运行：

```
./caddy start
```

## 参考资料

https://github.com/greenpau/caddy-trace