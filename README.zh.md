# Coffee

[English](README.md)

通过 ssh 隧道，将远程网络的 http(s) 服务中转成本地 http 服务。

## 使用场景

`coffee` 用来将远程网络 http(s) 服务中转成本地 http，因此您可以开始本
地编码和调试，无需在本地运行所有依赖的远程服务，因为这些服务很难甚至不
可能在本地运行：

1. 远程服务有自己的服务依赖或者中间件依赖。
2. 需要消耗极大的资源（cpu,内存）才能将他们全部在本地运行起来。
3. 需要耗费极大的精力来确保它们持续可用（试想依赖服务需要更新新的版本）。

总之一句话：你仅仅想这些服务直接稳定的被随时使用，这样你自己才能专注在
你的服务开发调试等工作上。

## 运行效果

- 帮助文档

如果不指定端口，`coffee` 将使用默认 `2046` 端口提供本地中继服务。他也
默认打印所有的请求响应日志到标准日志.

```bash
$ ./coffee --help
Usage of coffee:
  -port int
        use another serving port (default 2046)
  -trace
        trace every http roundtrip object (default true)
```

- 运行观察

启动 `coffee` 后，他将开始中继您配置的 http(s) 到远程网络。

```bash
$ ./coffee
coffee 🍵🍵🍵

The following remote network http(s) are relayed as local http:

http://www.vulnweb.com.127.0.0.1.nip.io:2046
http://rest.vulnweb.com.127.0.0.1.nip.io:2046
http://www.gnu.org.127.0.0.1.nip.io:2046
http://kernel.org.127.0.0.1.nip.io:2046
http://go.dev.127.0.0.1.nip.io:2046

🍵 2023/01/11 10:24:18 #1
HEAD / HTTP/1.1
Host: www.gnu.org
Accept: */*
User-Agent: curl/7.85.0


🍵 2023/01/11 10:24:18 establishing tunnel connection...
🍵 2023/01/11 10:24:20 #1
HTTP/1.1 200 OK
Connection: close
Accept-Ranges: bytes
Access-Control-Allow-Origin: (null)
Cache-Control: max-age=0
Content-Language: en
Content-Location: home.html
Content-Type: text/html
Date: Wed, 11 Jan 2023 02:24:19 GMT
Expires: Wed, 11 Jan 2023 02:24:19 GMT
Server: Apache/2.4.29
Strict-Transport-Security: max-age=63072000
Tcn: choice
Vary: negotiate,accept-language,Accept-Encoding
X-Content-Type-Options: nosniff
X-Frame-Options: sameorigin
```

- curl 请求将被中继到远程网络

以下是一个 `www.gnu.org` 本地中继后的情况。

```bash
# curl -I http://www.gnu.org.127.0.0.1.nip.io:2046
HTTP/1.1 200 OK
Accept-Ranges: bytes
Access-Control-Allow-Origin: (null)
Cache-Control: max-age=0
Content-Language: en
Content-Location: home.html
Content-Type: text/html
Date: Wed, 11 Jan 2023 02:24:19 GMT
Expires: Wed, 11 Jan 2023 02:24:19 GMT
Server: Apache/2.4.29
Strict-Transport-Security: max-age=63072000
Tcn: choice
Vary: negotiate,accept-language,Accept-Encoding
X-Content-Type-Options: nosniff
X-Frame-Options: sameorigin
```

## 配置我的环境 `coffee`

```json
{
  "wild": ".127.0.0.1.nip.io",
  "urls": [
    "http://www.vulnweb.com",
    "http://rest.vulnweb.com",
    "https://www.gnu.org",
    "https://kernel.org"
  ],
  "hops": [
    {
      "host": "127.0.0.1",
      "port": "22",
      "user": "ENV:USER",
      "pass": "ENV:PASS1",
      "key": "FILE:~/.ssh/id_rsa"
    },
    {
      "host": "localhost",
      "port": "22",
      "user": "ENV:USER",
      "pass": "ENV:PASS2",
      "key": "FILE:~/.ssh/id_rsa"
    }
  ]
}
```

### 编辑 `coffee.json`

1. 用你的远程网络地址替换 `urls` 数组
2. 配置 `coffee` 怎样一跳一跳的到达你的远程网络 (为保护你的隐私，`user`, `pass`, `key` 可选的支持环境变量或者文件读入)。
3. `wild` 可选的支持替换成任何可以解析到 `127.0.0.1` 的通配符域名后缀。

### 编译

编译最终的单二进制可执行文件，祝你使用 `coffee` 有一个愉快的过程！

```bash
make
```

## 技术参考引用

- [The Secure Shell (SSH) Protocol Architecture](https://www.rfc-editor.org/rfc/rfc4251)
- [The Secure Shell (SSH) Authentication Protocol](https://www.rfc-editor.org/rfc/rfc4252)
- [The Secure Shell (SSH) Transport Layer Protocol](https://www.rfc-editor.org/rfc/rfc4253)
- [The Secure Shell (SSH) Connection Protocol](https://www.rfc-editor.org/rfc/rfc4254)
- [The Secure Shell (SSH) Protocol Assigned Numbers](https://www.rfc-editor.org/rfc/rfc4250)
- [Hypertext Transfer Protocol -- HTTP/1.1](https://www.rfc-editor.org/rfc/rfc2616)
- [SSH Protocol by Go](https://pkg.go.dev/golang.org/x/crypto/ssh)
